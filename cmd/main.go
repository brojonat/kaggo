package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

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

func confirm(prompt string, r *bufio.Reader) (bool, error) {
	fmt.Println(prompt)
	input, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}
	txt := strings.TrimSpace(input)
	return slices.Contains([]string{"y", "yes", "ye"}, strings.ToLower(txt)), nil
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "admin",
				Usage: "Administrative commands (initiating workflows, etc.)",
				Subcommands: []*cli.Command{
					{
						Name:  "users",
						Usage: "Administrative user commands",
						Subcommands: []*cli.Command{
							{
								Name:  "add",
								Usage: "Add user",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "endpoint",
										Aliases:  []string{"end", "e"},
										Required: true,
										Usage:    "Kaggo server endpoint",
									},
									&cli.StringFlag{
										Name:     "email",
										Required: true,
										Usage:    "User's email",
									},
								},
								Action: func(ctx *cli.Context) error {
									return add_user(ctx)
								},
							},
							{
								Name:  "delete",
								Usage: "Delete user",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "endpoint",
										Aliases:  []string{"end", "e"},
										Required: true,
										Usage:    "Kaggo server endpoint",
									},
									&cli.StringFlag{
										Name:     "email",
										Required: true,
										Usage:    "User's email",
									},
								},
								Action: func(ctx *cli.Context) error {
									return delete_user(ctx)
								},
							},
							{
								Name:  "grant-metric",
								Usage: "Grant a metric to a user",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "endpoint",
										Aliases:  []string{"end", "e"},
										Required: true,
										Usage:    "Kaggo server endpoint",
									},
									&cli.StringFlag{
										Name:     "email",
										Required: true,
										Usage:    "User's email",
									},
									&cli.StringFlag{
										Name:     "request-kind",
										Aliases:  []string{"rk", "r"},
										Required: true,
										Usage:    "Metric request kind to grant",
									},
									&cli.StringFlag{
										Name:    "id",
										Aliases: []string{"i"},
										Usage:   "Metric identifier to grant",
									},
									&cli.BoolFlag{
										Name:    "all-ids",
										Aliases: []string{"all", "a"},
										Value:   false,
										Usage:   "Grant ALL metrics of request-kind to the user",
									},
								},
								Action: func(ctx *cli.Context) error {
									return grant_metric(ctx)
								},
							},
							{
								Name:  "remove-metric",
								Usage: "Remove a metric from a user",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "endpoint",
										Aliases:  []string{"end", "e"},
										Required: true,
										Usage:    "Kaggo server endpoint",
									},
									&cli.StringFlag{
										Name:     "email",
										Required: true,
										Usage:    "User's email",
									},
									&cli.StringFlag{
										Name:     "request-kind",
										Aliases:  []string{"rk", "r"},
										Required: true,
										Usage:    "Metric request kind to grant",
									},
									&cli.StringFlag{
										Name:    "id",
										Aliases: []string{"i"},
										Usage:   "Metric identifier to grant",
									},
									&cli.BoolFlag{
										Name:    "all-ids",
										Aliases: []string{"all", "a"},
										Value:   false,
										Usage:   "Grant ALL metrics of request-kind to the user",
									},
								},
								Action: func(ctx *cli.Context) error {
									return remove_metric(ctx)
								},
							},
						},
					},
					{
						Name:  "tinker",
						Usage: "testing sandbox/playground",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "test",
								Aliases: []string{"t"},
								Usage:   "test if passed to subcommand",
							},
						},
						Subcommands: []*cli.Command{
							{
								Name: "inner",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "test2", Aliases: []string{"t2"}, Usage: "testing"},
								},
								Action: func(ctx *cli.Context) error {
									return tinker(ctx)
								},
							},
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
								Usage:    "Request kind to perform",
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
						Name:  "delete-all-schedules",
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
							return delete_all_schedules(ctx)
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
						Usage: "Run the HTTP server on the specified port",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "listen-port",
								Aliases: []string{"port", "p"},
								Value:   os.Getenv("SERVER_PORT"),
								Usage:   "Port to listen on",
							},
							&cli.StringFlag{
								Name:    "database",
								Aliases: []string{"db", "d"},
								Value:   os.Getenv("DATABASE_URL"),
								Usage:   "Database endpoint",
							},
							&cli.StringFlag{
								Name:    "temporal-host",
								Aliases: []string{"th", "t"},
								Value:   os.Getenv("TEMPORAL_HOST"),
								Usage:   "Temporal endpoint",
							},
							&cli.IntFlag{
								Name:    "log-level",
								Aliases: []string{"ll", "l"},
								Usage:   "Logging level for the slog.Logger. Default is 0 (INFO), use -4 for DEBUG",
								Value:   0,
							},
						},
						Action: func(ctx *cli.Context) error {
							return serve_http(ctx)
						},
					},
					{
						Name:  "worker",
						Usage: "Run the Temporal worker",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "temporal-host",
								Aliases: []string{"th", "t"},
								Value:   os.Getenv("TEMPORAL_HOST"),
								Usage:   "Temporal endpoint",
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

func tinker(ctx *cli.Context) error {
	fmt.Println(ctx.String("test"), ctx.String("test2"))
	return nil
}
