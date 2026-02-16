package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/knr1997/assets-management-apiserver/internal/api/responses"
	"github.com/knr1997/assets-management-apiserver/internal/store"
)

type userKey string

const userCtx userKey = "user"

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

type UpdateUserPayload struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var payload UpdateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Username != nil {
		user.Username = *payload.Username
	}
	if payload.Email != nil {
		user.Email = *payload.Email
	}

	ctx := r.Context()

	if err := app.updateUser(ctx, user); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateUser(ctx context.Context, user *store.User) error {
	if err := app.store.Users.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (app *application) meDetailsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	response := responses.NewUserResponse(user)

	app.jsonResponse(w, http.StatusOK, response)
}

func (app *application) getAllUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := app.store.Users.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := responses.NewUsersResponse(users)

	// app.jsonResponse(w, http.StatusOK, resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
