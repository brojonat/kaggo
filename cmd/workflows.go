package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/api"
	"github.com/urfave/cli/v2"
)

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
