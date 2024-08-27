package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/brojonat/server-tools/stools"
)

func handleGetUsers(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emails := r.URL.Query()["email"]
		if len(emails) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply email(s)"))
			return
		}
		res, err := q.GetUsers(r.Context(), emails)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		if len(res) == 0 {
			writeEmptyResultError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func handleGetUserMetrics(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			writeBadRequestError(w, fmt.Errorf("must supply email(s)"))
			return
		}
		res, err := q.GetUserMetrics(r.Context(), email)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		if len(res) == 0 {
			writeEmptyResultError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func handleAddUser(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body api.CreateUserPayload
		err := stools.DecodeJSONBody(r, &body)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		if err = q.InsertUser(r.Context(), dbgen.InsertUserParams{Email: body.Email, Data: body.Data}); err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}

func handleDeleteUser(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emails := r.URL.Query()["email"]
		if len(emails) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply email(s)"))
			return
		}
		if err := q.DeleteUsers(r.Context(), emails); err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}

// switch over the set of allowed operations (add or remove metric). The API
// package will export supported operations.
func handleUserMetricOperation(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body api.UserMetricOperationPayload
		err := stools.DecodeJSONBody(r, &body)

		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		switch body.OpKind {
		case api.UserMetricOpKindAdd:
			p := dbgen.GrantMetricToUserParams{
				Email:       body.Email,
				RequestKind: body.RequestKind,
				ID:          body.ID,
			}
			if err = q.GrantMetricToUser(r.Context(), p); err != nil {
				if stools.IsPGError(err, stools.PGErrorForeignKeyViolation) {
					writeBadRequestError(w, fmt.Errorf("unable to grant; be sure user (%s) and metric (%s, %s) exists", body.Email, body.RequestKind, body.ID))
					return
				}
				writeInternalError(l, w, err)
				return
			}
			writeOK(w)
			return

		case api.UserMetricOpKindAddGroup:
			p := dbgen.GrantMetricGroupToUserParams{
				Email:       body.Email,
				RequestKind: body.RequestKind,
			}
			if err = q.GrantMetricGroupToUser(r.Context(), p); err != nil {
				if stools.IsPGError(err, stools.PGErrorForeignKeyViolation) ||
					stools.IsPGError(err, stools.PGErrorUniqueViolation) {
					writeBadRequestError(w, fmt.Errorf("unable to grant; be sure user (%s) and metric (%s) exists", body.Email, body.RequestKind))
					return
				}
				writeInternalError(l, w, err)
				return
			}
			writeOK(w)
			return

		case api.UserMetricOpKindRemove:
			p := dbgen.RemoveMetricFromUserParams{
				Email:       body.Email,
				RequestKind: body.RequestKind,
				ID:          body.ID,
			}
			if err = q.RemoveMetricFromUser(r.Context(), p); err != nil {
				writeInternalError(l, w, err)
				return
			}
			writeOK(w)
			return

		case api.UserMetricOpKindRemoveGroup:
			p := dbgen.RemoveMetricGroupFromUserParams{
				Email:       body.Email,
				RequestKind: body.RequestKind,
			}
			if err = q.RemoveMetricGroupFromUser(r.Context(), p); err != nil {
				writeInternalError(l, w, err)
				return
			}
			writeOK(w)
			return

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported op_kind: %s", body.OpKind))
			return
		}
	}
}
