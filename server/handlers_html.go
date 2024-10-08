package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/jackc/pgx/v5/pgtype"
	"go.temporal.io/sdk/client"
)

const (
	PlotKindScheduleCount    string = "schedule-count"
	PlotKindScheduleTimeline string = "schedule-timeline"
	PlotKindUserPulse        string = "user-pulse"
)

var plotTmpl *template.Template
var d3Tmpl *template.Template

type templateData struct {
	Endpoint                 string
	PlotKind                 string
	CDN                      string
	LocalStorageAuthTokenKey string
	DataPath                 string
	ID                       string
}

func handleGetPlotData(l *slog.Logger, q *dbgen.Queries, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pk := r.URL.Query().Get("plot_kind")
		switch pk {
		case PlotKindScheduleCount:

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

		case PlotKindScheduleTimeline:
			ss, err := tc.ScheduleClient().List(r.Context(), client.ScheduleListOptions{})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

			type schedule struct {
				ID      string
				Nexts   []time.Time
				Jitter  time.Duration
				Comment string
			}
			schedules := []schedule{}

			for {
				if !ss.HasNext() {
					break
				}
				s, err := ss.Next()
				if err != nil {
					break
				}
				parts := strings.Split(s.ID, " ")
				id := fmt.Sprintf("%s %s", parts[0], parts[1])
				nexts := getNextActionTimesNoJitter(s.Spec, 3)
				i := schedule{
					ID:      id,
					Nexts:   nexts,
					Jitter:  s.Spec.Jitter,
					Comment: s.Spec.Calendars[0].Comment,
				}
				schedules = append(schedules, i)
			}

			// construct the traces
			traces := []types.Trace{}
			for _, s := range schedules {

				xs := []time.Time{}
				xls := []time.Time{}
				ys := []string{}

				for _, n := range s.Nexts {
					if n.Before(time.Now()) {
						continue
					}
					xs = append(xs, n)
					xls = append(xls, n.Add(s.Jitter))
					ys = append(ys, s.ID)
				}

				xs = xs[0:10]
				xls = xls[0:10]
				ys = ys[0:10]

				traces = append(traces, &grob.Scatter{
					Marker:      &grob.ScatterMarker{Size: types.ArrayOKValue(types.N(20))},
					X:           types.DataArray(xs),
					Y:           types.DataArray(ys),
					Name:        types.S(s.ID),
					Legendgroup: types.S(s.ID),
					Mode:        grob.ScatterModeMarkers,
				})
				traces = append(traces, &grob.Scatter{
					Marker: &grob.ScatterMarker{
						Color:  types.ArrayOKValue(types.UseColor(types.C("black"))),
						Symbol: types.ArrayOKValue(grob.ScatterMarkerSymbolCircleXOpen), Size: types.ArrayOKValue(types.N(10))},
					X:           types.DataArray(xls),
					Y:           types.DataArray(ys),
					Name:        types.S(s.ID),
					Legendgroup: types.S(s.ID),
					Showlegend:  types.False,
					Mode:        grob.ScatterModeMarkers,
				})

			}

			fig := &grob.Fig{
				Data: traces,
				Layout: &grob.Layout{
					PlotBgcolor: "whitesmoke",
					Height:      types.N(20000),
					Margin:      &grob.LayoutMargin{Autoexpand: types.True, L: types.N(500)},
					Xaxis:       &grob.LayoutXaxis{Title: &grob.LayoutXaxisTitle{Text: types.S("Time")}},
					Yaxis:       &grob.LayoutYaxis{Title: &grob.LayoutYaxisTitle{Text: types.S("Schedule")}},
					Title: &grob.LayoutTitle{
						Text: types.S("When will schedules be running?"),
					},
					Legend: &grob.LayoutLegend{},
				},
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fig)

		case PlotKindUserPulse:
			id := r.URL.Query().Get("id")
			user_metrics, err := q.GetRedditUserMetricsByIDsBucket15Min(
				r.Context(), dbgen.GetRedditUserMetricsByIDsBucket15MinParams{
					Ids:     []string{id},
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user_metrics)
			return

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported plot_kind %s", pk))
			return
		}
	}
}

func handleGetPlots(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pk := r.URL.Query().Get("plot_kind")
		id := r.URL.Query().Get("id")
		q := url.Values{}

		switch pk {
		case PlotKindScheduleCount, PlotKindScheduleTimeline:
			// Plots schedule related information. These plots use an internal
			// handler specifically designed to interface with the Temporal
			// client and contruct the relevant data.
			q.Add("plot_kind", pk)
			data := templateData{
				Endpoint:                 os.Getenv("KAGGO_ENDPOINT"),
				DataPath:                 "/plot-data?" + q.Encode(),
				CDN:                      os.Getenv("PLOTLY_CDN"),
				LocalStorageAuthTokenKey: os.Getenv("LOCAL_STORAGE_AUTH_TOKEN_KEY"),
			}
			w.WriteHeader(http.StatusOK)
			err := plotTmpl.Execute(w, data)
			if err != nil {
				l.Error("Error rendering template", "error", err)
				writeInternalError(l, w, err)
			}
		case PlotKindUserPulse:
			// Plot shows the user pulse metrics. This plot fetches the data
			// from the backend timeseries handlers directly.
			if id == "" {
				writeBadRequestError(w, fmt.Errorf("must supply request_kind and id"))
				return
			}
			q.Add("plot_kind", pk)
			q.Add("id", id)
			data := templateData{
				Endpoint:                 os.Getenv("KAGGO_ENDPOINT"),
				PlotKind:                 pk,
				CDN:                      os.Getenv("PLOTLY_CDN"),
				LocalStorageAuthTokenKey: os.Getenv("LOCAL_STORAGE_AUTH_TOKEN_KEY"),
				ID:                       id,
			}
			w.WriteHeader(http.StatusOK)
			err := d3Tmpl.Execute(w, data)
			if err != nil {
				l.Error("Error rendering template", "error", err)
				writeInternalError(l, w, err)
			}

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported plot_kind %s", pk))
			return
		}

	}
}

// This is a broken but useful helper function. This assumes that the supplied
// ScheduleSpec has one calendar with one schedule range, no skips, and that
// every day has the same schedule. This is sufficient for now but ultimately
// will need to be improved to handle more complicated schedules. Ideally
// temporal would expose this functionality but alas.
func getNextActionTimesNoJitter(s *client.ScheduleSpec, ndays int) []time.Time {

	seconds := []int{}
	minutes := []int{}
	hours := []int{}

	for t := s.Calendars[0].Second[0].Start; t <= s.Calendars[0].Second[0].End; t += s.Calendars[0].Second[0].Step {
		seconds = append(seconds, t)
	}
	for t := s.Calendars[0].Minute[0].Start; t <= s.Calendars[0].Minute[0].End; t += s.Calendars[0].Minute[0].Step {
		minutes = append(minutes, t)
	}
	for t := s.Calendars[0].Hour[0].Start; t <= s.Calendars[0].Hour[0].End; t += s.Calendars[0].Hour[0].Step {
		hours = append(hours, t)
	}

	// for the next 5 days, get all the scheduled runs
	now := time.Now()
	nexts := []time.Time{}
	for nday := range ndays {
		for _, h := range hours {
			for _, m := range minutes {
				for _, sec := range seconds {
					date := now.AddDate(0, 0, nday)
					itert := time.Date(date.Year(), date.Month(), date.Day(), h, m, sec, 0, date.Location())
					nexts = append(nexts, itert)
				}
			}
		}
	}
	return nexts
}
