package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brojonat/kaggo/server"
	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/worker"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
)

func getDefaultLogger(lvl slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.Function = ""
					source.File = filepath.Base(source.File)
				}
			}
			return a
		},
	}))
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "admin",
				Usage: "Administrative commands (initiating workflows, etc.)",
				Subcommands: []*cli.Command{
					{
						Name:  "tinker",
						Usage: "playground",
						Flags: []cli.Flag{},
						Action: func(ctx *cli.Context) error {
							return tinker(ctx)
						},
					},
					{
						Name:  "run-metadata-workflow",
						Usage: "Send a POST request to initiate a metadata workflow",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "endpoint",
								Aliases:  []string{"end", "e"},
								Required: true,
								Usage:    "Kaggo server endpoint",
							},
							&cli.StringFlag{
								Name:     "request-kind",
								Aliases:  []string{"rk", "r"},
								Required: true,
								Usage:    "Request kind to perform.",
							},
							&cli.StringFlag{
								Name:     "id",
								Aliases:  []string{"i"},
								Required: true,
								Usage:    "Resource ID for the request",
							},
						},
						Action: func(ctx *cli.Context) error {
							return run_metadata_wf(ctx)
						},
					},
					{
						Name:  "dump-schedules",
						Usage: "Dump schedules to file",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "endpoint",
								Aliases:  []string{"end", "e"},
								Required: true,
								Usage:    "Kaggo server endpoint",
							},
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Required: true,
								Usage:    "Output file location",
							},
						},
						Action: func(ctx *cli.Context) error {
							return dump_schedules(ctx)
						},
					},
					{
						Name:  "delete-schedules",
						Usage: "Delete all schedules. Be sure to dump a backup first!",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "endpoint",
								Aliases:  []string{"end", "e"},
								Required: true,
								Usage:    "Kaggo server endpoint",
							},
						},
						Action: func(ctx *cli.Context) error {
							return delete_schedules(ctx)
						},
					},
					{
						Name:  "load-schedules",
						Usage: "Load schedules from file",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "endpoint",
								Aliases:  []string{"end", "e"},
								Required: true,
								Usage:    "Kaggo server endpoint",
							},
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Required: true,
								Usage:    "Input file location",
							},
						},
						Action: func(ctx *cli.Context) error {
							return load_schedules(ctx)
						},
					},
				},
			},
			{
				Name:  "run",
				Usage: "Commands for running various components (server, workers, etc.)",
				Subcommands: []*cli.Command{
					{
						Name:  "http-server",
						Usage: "Run the HTTP server on the specified port.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "listen-port",
								Aliases: []string{"port", "p"},
								Value:   os.Getenv("SERVER_PORT"),
								Usage:   "Port to listen on.",
							},
							&cli.StringFlag{
								Name:    "database",
								Aliases: []string{"db", "d"},
								Value:   os.Getenv("DATABASE_URL"),
								Usage:   "Database endpoint.",
							},
							&cli.StringFlag{
								Name:    "temporal-host",
								Aliases: []string{"th", "t"},
								Value:   os.Getenv("TEMPORAL_HOST"),
								Usage:   "Temporal endpoint.",
							},
							&cli.IntFlag{
								Name:    "log-level",
								Aliases: []string{"ll", "l"},
								Usage:   "Logging level for the slog.Logger. Default is 0 (INFO), use -4 for DEBUG.",
								Value:   0,
							},
						},
						Action: func(ctx *cli.Context) error {
							return serve_http(ctx)
						},
					},
					{
						Name:  "worker",
						Usage: "Run the Temporal worker.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "temporal-host",
								Aliases: []string{"th", "t"},
								Value:   os.Getenv("TEMPORAL_HOST"),
								Usage:   "Temporal endpoint.",
							},
						},
						Action: func(ctx *cli.Context) error {
							return run_worker(ctx)
						},
					},
				},
			},
		}}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error running command: %s\n", err.Error())
		os.Exit(1)
	}
}

func serve_http(ctx *cli.Context) error {
	// internal init
	logger := getDefaultLogger(slog.Level(ctx.Int("log-level")))

	return server.RunHTTPServer(
		ctx.Context,
		ctx.String("listen-port"),
		logger,
		ctx.String("database"),
		ctx.String("temporal-host"),
	)
}

func run_worker(ctx *cli.Context) error {
	logger := getDefaultLogger(slog.Level(ctx.Int("log-level")))
	thp := ctx.String("temporal-host")
	return worker.RunWorker(ctx.Context, logger, thp)
}

func run_metadata_wf(ctx *cli.Context) error {
	p := api.GenericScheduleRequestPayload{
		RequestKind: ctx.String("request-kind"),
		ID:          ctx.String("id"),
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not serialize payload: %w", err)
	}
	r, err := http.NewRequest(
		http.MethodPost,
		ctx.String("endpoint")+"/metadata/run-workflow",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s", res.Status)
	}
	return nil
}

func dump_schedules(ctx *cli.Context) error {
	r, err := http.NewRequest(
		http.MethodGet,
		ctx.String("endpoint")+"/schedule",
		nil,
	)
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(ctx.String("file"), b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func delete_schedules(ctx *cli.Context) error {

	r, err := http.NewRequest(
		http.MethodGet,
		ctx.String("endpoint")+"/schedule",
		nil,
	)
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var schedules []struct {
		ID   string              `json:"ID"`
		Spec client.ScheduleSpec `json:"Spec"`
	}
	err = json.Unmarshal(b, &schedules)
	if err != nil {
		return err
	}

	for i, sched := range schedules {
		r, err := http.NewRequest(
			http.MethodDelete,
			ctx.String("endpoint")+"/schedule",
			nil,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error making request to schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		q := r.URL.Query()
		q.Add("schedule_id", sched.ID)
		r.URL.RawQuery = q.Encode()

		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error deleting schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		defer res.Body.Close()
		b, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading response for schedule delete %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		var rbody api.DefaultJSONResponse
		err = json.Unmarshal(b, &rbody)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing \"%s\" response for schedule delete %d (%s): %s\n", res.Status, i, sched.ID, err.Error())
			continue
		}
		if res.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "%s response deleting schedule %d (%s): %s\n", res.Status, i, sched.ID, rbody.Error)
			continue
		}
	}
	return nil
}

func load_schedules(ctx *cli.Context) error {
	b, err := os.ReadFile(ctx.String("file"))
	if err != nil {
		return err
	}
	var body []struct {
		ID   string              `json:"ID"`
		Spec client.ScheduleSpec `json:"Spec"`
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	}
	for i, sched := range body {
		parts := strings.Split(sched.ID, " ")
		payload := api.GenericScheduleRequestPayload{
			RequestKind: parts[0],
			ID:          parts[1],
			Schedule:    sched.Spec,
		}
		b, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating payload for schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		r, err := http.NewRequest(
			http.MethodPost,
			ctx.String("endpoint")+"/schedule",
			bytes.NewReader(b),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error making request to schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error uploading schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		defer res.Body.Close()
		b, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading response for schedule upload %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		var rbody api.DefaultJSONResponse
		err = json.Unmarshal(b, &rbody)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing \"%s\" response for schedule upload %d (%s): %s\n", res.Status, i, sched.ID, err.Error())
			continue
		}
		if res.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "%s response uploading schedule %d (%s): %s\n", res.Status, i, sched.ID, rbody.Error)
			continue
		}
	}
	return nil
}

func tinker(ctx *cli.Context) error {
	td, _ := time.ParseDuration("5s")
	s := client.ScheduleSpec{Jitter: td}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
	return nil
}
