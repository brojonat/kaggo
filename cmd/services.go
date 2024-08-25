package main

import (
	"log/slog"

	"github.com/brojonat/kaggo/server"
	"github.com/brojonat/kaggo/worker"
	"github.com/urfave/cli/v2"
)

func serve_http(ctx *cli.Context) error {
	// internal init
	logger := getDefaultLogger(slog.Level(ctx.Int("log-level")))
	pms := server.GetDefaultPromMetrics()

	return server.RunHTTPServer(
		ctx.Context,
		ctx.String("listen-port"),
		logger,
		ctx.String("database"),
		ctx.String("temporal-host"),
		pms,
	)
}

func run_worker(ctx *cli.Context) error {
	logger := getDefaultLogger(slog.Level(ctx.Int("log-level")))
	thp := ctx.String("temporal-host")
	return worker.RunWorker(ctx.Context, logger, thp)
}
