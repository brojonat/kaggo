package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/brojonat/kaggo/server/api"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
)

func dump_schedules(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodGet, ctx.String("endpoint")+"/schedule", nil)
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

func delete_all_schedules(ctx *cli.Context) error {
	c, err := confirm("Delete all schedules!?", bufio.NewReader(os.Stdin))
	if !c || err != nil {
		return fmt.Errorf("confirmation failed, aborting")
	}
	r, err := http.NewRequest(http.MethodGet, ctx.String("endpoint")+"/schedule", nil)
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
		r, err := http.NewRequest(http.MethodDelete, ctx.String("endpoint")+"/schedule", nil)
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

		// NOTE: don't block for the metadata operation. We don't want that to
		// impact the schedule creation here because this is typically used to
		// bulk (re)create schedules. If you want to bulk create schedules that
		// need their metadata fetched, then you'll need to thread that update
		// the CLI to accept that as an additional flag and pass it here.
		q := r.URL.Query()
		q.Add("skip-metadata", "true")
		r.URL.RawQuery = q.Encode()
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

func create_schedule(ctx *cli.Context) error {
	rk := ctx.String("request-kind")
	id := ctx.String("id")
	sched := kt.GetDefaultScheduleSpec(rk, id)
	payload := api.GenericScheduleRequestPayload{
		RequestKind: rk,
		ID:          id,
		Schedule:    sched,
		Monitor:     ctx.Bool("monitor"),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/schedule", bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %d: %s", res.StatusCode, string(b))
	}
	return nil
}

func delete_schedule(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodDelete, ctx.String("endpoint")+"/schedule", nil)
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	q := r.URL.Query()
	q.Add("schedule_id", ctx.String("schedule_id"))
	r.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %d: %s", res.StatusCode, string(b))
	}
	return nil
}

// This will read a schedule dump file, iterate over the schedules, and for any
// schedules matching the supplied request kind, it will delete the existing
// schedule and immediately re-create it. This is useful for situations such as
// after we've changed the default schedule configuration and we'd like to
// replace existing schedules with schedules that have the new configuration.
func reupload_schedules(ctx *cli.Context) error {
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

	// iterate over all the schedules in the file
	for i, sched := range body {

		parts := strings.Split(sched.ID, " ")
		rk := parts[0]
		id := parts[1]

		// no-op if the schedule doesn't match the supplied kind
		if rk != ctx.String("request-kind") {
			continue
		}

		// first delete
		r, err := http.NewRequest(http.MethodDelete, ctx.String("endpoint")+"/schedule", nil)
		if err != nil {
			return err
		}
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
		q := r.URL.Query()
		q.Add("schedule_id", sched.ID)
		r.URL.RawQuery = q.Encode()
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("bad response from server: %d: %s", res.StatusCode, string(b))
		}

		// now re-upload using a (new) default schedule spec
		payload := api.GenericScheduleRequestPayload{
			RequestKind: rk,
			ID:          id,
			Schedule:    kt.GetDefaultScheduleSpec(rk, id),
		}
		b, err = json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating payload for schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}
		r, err = http.NewRequest(
			http.MethodPost,
			ctx.String("endpoint")+"/schedule",
			bytes.NewReader(b),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error making request to schedule %d (%s): %s\n", i, sched.ID, err.Error())
			continue
		}

		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))

		// NOTE: don't block for the metadata operation. We don't want that to
		// impact the schedule creation here because this is typically used to
		// bulk (re)create schedules. If you want to bulk create schedules that
		// need their metadata fetched, then you'll need to update the CLI to
		// accept that as an additional flag and pass it here.
		q = r.URL.Query()
		q.Add("skip-metadata", "true")
		r.URL.RawQuery = q.Encode()
		res, err = http.DefaultClient.Do(r)
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
