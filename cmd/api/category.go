package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/knr1997/assets-management-apiserver/internal/store"
)

type categoryKey string

const categoryCtx categoryKey = "category"

type CreateCategoryPayload struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCategoryPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Category{
		Name:        payload.Name,
		Description: payload.Description,
	}

	ctx := r.Context()

	if err := app.store.Category.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdateCategoryPayload struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description"`
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := getCategoryFromCtx(r)

	var payload UpdateCategoryPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Name != nil {
		category.Name = *payload.Name
	}
	if payload.Description != nil {
		category.Description = *payload.Description
	}

	ctx := r.Context()

	if err := app.updateCategory(ctx, category); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, category); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := getCategoryFromCtx(r)

	response := ToCategoryResponse(category)

	if err := app.jsonResponse(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CategoryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ToCategoryResponse(a *store.Category) CategoryResponse {
	return CategoryResponse{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
	}
}

func ToCategoryResponseList(categories []store.Category) []CategoryResponse {
	responses := make([]CategoryResponse, 0, len(categories))

	for i := range categories {
		responses = append(responses, ToCategoryResponse(&categories[i]))
	}

	return responses
}

func (app *application) getAllCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := app.store.Category.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := ToCategoryResponseList(categories)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "categoryID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Category.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) categoryContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "categoryID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		category, err := app.store.Category.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, categoryCtx, category)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCategoryFromCtx(r *http.Request) *store.Category {
	category, _ := r.Context().Value(categoryCtx).(*store.Category)
	return category
}

func (app *application) updateCategory(ctx context.Context, category *store.Category) error {
	if err := app.store.Category.Update(ctx, category); err != nil {
		return err
	}

	return nil
}
