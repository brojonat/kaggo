package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/api"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/urfave/cli/v2"
)

func add_listener_subscription(ctx *cli.Context) error {
	rk := ctx.String("request-kind")
	id := ctx.String("id")

	if rk != kt.RequestKindYouTubeChannel {
		return fmt.Errorf("unsupported request kind %s", rk)
	}

	var b []byte
	var err error

	p := api.AddListenerSubPayload{
		RequestKind: rk,
		ID:          id,
	}
	b, err = json.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not serialize subscription payload: %w", err)
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/add-listener-sub", bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s", res.Status)
	}
	return nil
}
