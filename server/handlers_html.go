package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"github.com/brojonat/kaggo/server/db/dbgen"
)

var plotTmpl *template.Template

type templateData struct {
	CDN                      string
	ContentURL               string
	LocalStorageAuthTokenKey string
	Endpoint                 string
}

func htmlWriteTemplate(
	w http.ResponseWriter,
	t *template.Template,
	statusCode int,
	td any,
) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html")
	return t.Execute(w, td)
}

func handlePlotData(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("id") {
		case "workflow-count-projection":
			fig := &grob.Fig{
				Data: []types.Trace{
					&grob.Bar{
						X: types.DataArray([]float64{1, 2, 3}),
						Y: types.DataArray([]float64{42, 64, 128}),
					},
				},
				Layout: &grob.Layout{
					PlotBgcolor: "chartreuse",
					Xaxis:       &grob.LayoutXaxis{Title: &grob.LayoutXaxisTitle{Text: types.S("Hi")}},
					Yaxis:       &grob.LayoutYaxis{Title: &grob.LayoutYaxisTitle{Text: types.S("Hi")}},
					Title: &grob.LayoutTitle{
						Text: types.S("How Many Workflows Will Be Running?"),
					},
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
		switch r.URL.Query().Get("id") {
		case "workflow-count-projection":
			data := templateData{
				CDN:                      os.Getenv("PLOTLY_CDN"),
				ContentURL:               os.Getenv("KAGGO_ENDPOINT") + "/plot-data?id=dummy",
				LocalStorageAuthTokenKey: os.Getenv("LOCAL_STORAGE_AUTH_TOKEN_KEY"),
				Endpoint:                 os.Getenv("KAGGO_ENDPOINT"),
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
