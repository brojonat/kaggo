package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/brojonat/kaggo/server/api"
	"github.com/urfave/cli/v2"
)

func tinker_wf(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/run-reddit-listener-wf", nil)
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

func initiate_youtube_listener(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/run-youtube-listener-wf", nil)
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

func run_metadata_wf(ctx *cli.Context) error {
	id := ctx.String("id")
	rk := ctx.String("request-kind")
	allIDs := ctx.Bool("all-ids")

	if id == "" && !allIDs {
		return fmt.Errorf("must supply an id or --all-ids")
	}

	var ids []string
	if allIDs {
		// get all ids for the request kind
		r, err := http.NewRequest(http.MethodGet, ctx.String("endpoint")+"/schedule", nil)
		if err != nil {
			return err
		}
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
		q := r.URL.Query()
		q.Add("request_kind", rk)
		r.URL.RawQuery = q.Encode()
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("bad response from server: %s", res.Status)
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		var body []struct {
			ID string
		}
		err = json.Unmarshal(b, &body)
		if err != nil {
			return fmt.Errorf("could not deserialize schedule response")
		}
		for _, s := range body {
			ids = append(ids, s.ID)
		}
	} else {
		// construct the ID as if it came back from the temporal server
		ids = []string{fmt.Sprintf("%s %s", rk, id)}
	}

	// range over ids and kick off the metadata workflow
	successCount := 0
	for _, id := range ids {
		id_parts := strings.Split(id, " ")
		p := api.GenericScheduleRequestPayload{
			RequestKind: ctx.String("request-kind"),
			ID:          id_parts[1],
		}
		body, err := json.Marshal(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not serialize payload: %s", err)
			continue
		}
		r, err := http.NewRequest(
			http.MethodPost,
			ctx.String("endpoint")+"/metadata/run-workflow",
			bytes.NewReader(body),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error constructing request for schedule %s: %s", id, err)
			continue
		}
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error doing request for schedule %s: %s", id, err)
			continue
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "bad response from server for schedule %s: %s", id, res.Status)
			continue
		}
		successCount++
	}
	if successCount != len(ids) {
		return fmt.Errorf("success count %d / %d", successCount, len(ids))
	}
	return nil
}
