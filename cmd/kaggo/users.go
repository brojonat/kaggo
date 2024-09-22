package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/api"
	"github.com/urfave/cli/v2"
)

func add_user(ctx *cli.Context) error {
	p := api.CreateUserPayload{
		Email: ctx.String("email"),
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not serialize payload: %w", err)
	}
	r, err := http.NewRequest(
		http.MethodPost,
		ctx.String("endpoint")+"/users",
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
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s: %s", res.Status, body)
	}
	return nil
}

func delete_user(ctx *cli.Context) error {
	r, err := http.NewRequest(
		http.MethodDelete,
		ctx.String("endpoint")+"/users",
		nil,
	)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	q.Add("email", ctx.String("email"))
	r.URL.RawQuery = q.Encode()
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s: %s", res.Status, body)
	}
	return nil
}

func grant_metric(ctx *cli.Context) error {
	var body []byte
	var err error
	if ctx.Bool("all-ids") {
		p := api.UserMetricOperationPayload{
			OpKind:      api.UserMetricOpKindAddGroup,
			Email:       ctx.String("email"),
			RequestKind: ctx.String("request-kind"),
		}
		body, err = json.Marshal(p)
	} else {
		p := api.UserMetricOperationPayload{
			OpKind:      api.UserMetricOpKindAdd,
			Email:       ctx.String("email"),
			RequestKind: ctx.String("request-kind"),
			ID:          ctx.String("id"),
		}
		body, err = json.Marshal(p)
	}

	if err != nil {
		return fmt.Errorf("could not serialize payload: %w", err)
	}
	r, err := http.NewRequest(
		http.MethodPost,
		ctx.String("endpoint")+"/users/metrics",
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
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s: %s", res.Status, body)
	}
	return nil
}

func remove_metric(ctx *cli.Context) error {
	var body []byte
	var err error
	if ctx.Bool("all-ids") {
		p := api.UserMetricOperationPayload{
			OpKind:      api.UserMetricOpKindRemoveGroup,
			Email:       ctx.String("email"),
			RequestKind: ctx.String("request-kind"),
		}
		body, err = json.Marshal(p)
	} else {
		p := api.UserMetricOperationPayload{
			OpKind:      api.UserMetricOpKindRemove,
			Email:       ctx.String("email"),
			RequestKind: ctx.String("request-kind"),
			ID:          ctx.String("id"),
		}
		body, err = json.Marshal(p)
	}
	if err != nil {
		return fmt.Errorf("could not serialize payload: %w", err)
	}
	r, err := http.NewRequest(
		http.MethodPost,
		ctx.String("endpoint")+"/users/metrics",
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
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s: %s", res.Status, body)
	}
	return nil
}
