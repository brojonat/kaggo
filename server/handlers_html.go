package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"go.temporal.io/sdk/client"
)

const PlotKindWorkflowCountProjection string = "workflow-count-projection"

var plotTmpl *template.Template

type templateData struct {
	Endpoint                 string
	PlotKind                 string
	CDN                      string
	LocalStorageAuthTokenKey string
}

func handlePlotData(l *slog.Logger, q *dbgen.Queries, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("id") {
		case PlotKindWorkflowCountProjection:

			ss, err := tc.ScheduleClient().List(r.Context(), client.ScheduleListOptions{})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

			type schedule struct {
				RequestKind string
				StartTS     float64
				EndTS       float64
			}
			data := []schedule{}

			for {
				if !ss.HasNext() {
					break
				}
				s, err := ss.Next()
				if err != nil {
					break
				}
				rk := strings.Split(s.ID, " ")[0]
				st := s.Spec.StartAt.Unix()
				et := s.Spec.EndAt.Unix()
				i := schedule{
					RequestKind: rk,
					StartTS:     float64(st),
					EndTS:       float64(et),
				}
				data = append(data, i)
			}

			truncateToDay := func(t time.Time) time.Time {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			}

			// Construct the abscissa by starting at the beginning of today and
			// adding nbins that are binsize wide rounded to the nearest minute.
			binsize := 24 * 3600
			nbins := 30
			abscissa := make([]float64, nbins)
			for i := 0; i < nbins; i++ {
				offset := time.Duration(time.Duration(i*binsize) * time.Second)
				abscissa[i] = float64(truncateToDay(time.Now()).Add(offset.Round(time.Minute)).Unix())
			}

			// Loop over the abscissa and for each bin, count the number of schedules that
			// will be running in that bin. This is partitioned by RequestKind.
			counts := map[string][]float64{}
			for i, x := range abscissa {
				binUpper := x + float64(binsize)
				for _, s := range data {
					// schedules with no start/end time defined will have the
					// zero value of the time, so handle those cases as well as
					// the naive check for bin containment
					startsBefore := bool(s.StartTS < x)
					doesNotEnd := time.Unix(int64(s.EndTS), 0).IsZero()
					runsIn := bool(s.StartTS <= x && s.EndTS > x)
					runsOut := bool(s.StartTS < binUpper && s.EndTS >= x)
					enclosed := bool(s.StartTS >= x && s.EndTS <= binUpper)
					if (startsBefore && doesNotEnd) || runsIn || runsOut || enclosed {
						if _, ok := counts[s.RequestKind]; !ok {
							counts[s.RequestKind] = make([]float64, nbins)
						}
						counts[s.RequestKind][i] += 1
					}
				}
			}

			// relabel x axis
			tss := []string{}
			for _, a := range abscissa {
				tss = append(tss, time.Unix(int64(a), 0).Format(time.RFC3339))
			}

			// construct the traces
			traces := []types.Trace{}
			for _, rk := range kt.GetSupportedRequestKinds() {
				traces = append(traces, &grob.Bar{
					X:    types.DataArray(tss),
					Y:    types.DataArray(counts[rk]),
					Name: types.S(rk),
				})
			}

			fig := &grob.Fig{
				Data: traces,
				Layout: &grob.Layout{
					PlotBgcolor: "whitesmoke",
					Xaxis:       &grob.LayoutXaxis{Title: &grob.LayoutXaxisTitle{Text: types.S("Time")}},
					Yaxis:       &grob.LayoutYaxis{Title: &grob.LayoutYaxisTitle{Text: types.S("Number of Schedules Running")}},
					Title: &grob.LayoutTitle{
						Text: types.S("How Many Schedules Will Be Running?"),
					},
					Legend: &grob.LayoutLegend{},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fig)
			return

		default:
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
	}
}

func handlePlot(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pk := r.URL.Query().Get("id")
		switch pk {
		case PlotKindWorkflowCountProjection:
			data := templateData{
				Endpoint:                 os.Getenv("KAGGO_ENDPOINT"),
				PlotKind:                 pk,
				CDN:                      os.Getenv("PLOTLY_CDN"),
				LocalStorageAuthTokenKey: os.Getenv("LOCAL_STORAGE_AUTH_TOKEN_KEY"),
			}
			w.WriteHeader(http.StatusOK)
			err := plotTmpl.Execute(w, data)
			if err != nil {
				l.Error("Error rendering template", "error", err)
			}
		default:
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}

	}
}
