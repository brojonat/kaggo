package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/brojonat/kaggo/server"
	"github.com/brojonat/kaggo/worker"
	"github.com/urfave/cli/v2"
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
				Name:  "tinker",
				Usage: "playground",
				Flags: []cli.Flag{},
				Action: func(ctx *cli.Context) error {
					return tinker(ctx)
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
	metrics := server.GetDefaultPromMetrics()

	return server.RunHTTPServer(
		ctx.Context,
		ctx.String("listen-port"),
		logger,
		ctx.String("database"),
		ctx.String("temporal-host"),
		metrics,
	)
}

func run_worker(ctx *cli.Context) error {
	logger := getDefaultLogger(slog.Level(ctx.Int("log-level")))
	thp := ctx.String("temporal-host")
	return worker.RunWorker(ctx.Context, logger, thp)
}

func tinker(ctx *cli.Context) error {
	fmt.Println("hi")
	return nil
}
