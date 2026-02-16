package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/knr1997/assets-management-apiserver/internal/api/requests"
	"github.com/knr1997/assets-management-apiserver/internal/api/responses"
	"github.com/knr1997/assets-management-apiserver/internal/store"
)

type assetKey string

const assetCtx assetKey = "asset"

type CreateAssetPayload struct {
	Name         string `json:"name" validate:"required,max=100"`
	SerialNumber string `json:"serialNumber" validate:"required,max=100"`
	Description  string `json:"description"`
	CategoryID   int64  `json:"categoryID" validate:"required"`
}

func (app *application) createAssetHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateAssetPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	asset := &store.Asset{
		Name:         payload.Name,
		SerialNumber: payload.SerialNumber,
		Description:  payload.Description,
		CategoryID:   payload.CategoryID,
	}

	ctx := r.Context()

	if err := app.store.Asset.Create(ctx, asset); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, asset); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdateAssetPayload struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description"`
}

func (app *application) updateAssetHandler(w http.ResponseWriter, r *http.Request) {
	asset := getAssetFromCtx(r)

	var payload UpdateAssetPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Name != nil {
		asset.Name = *payload.Name
	}
	if payload.Description != nil {
		asset.Description = *payload.Description
	}

	ctx := r.Context()

	if err := app.updateAsset(ctx, asset); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, asset); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getAssetHandler(w http.ResponseWriter, r *http.Request) {
	asset := getAssetFromCtx(r)

	response := responses.NewAssetResponse(asset)

	if err := app.jsonResponse(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getAllAssetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	assets, err := app.store.Asset.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := responses.NewAssetsResponse(assets)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) deleteAssetHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "assetID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Asset.Delete(ctx, id); err != nil {
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

func (app *application) checkoutAssetHandler(w http.ResponseWriter, r *http.Request) {
	asset := getAssetFromCtx(r)
	var payload requests.CheckoutAssetPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 🔒 Business rule
	if asset.Status != store.AssetAvailable {
		app.badRequestResponse(w, r, errors.New("asset is not available"))
		return
	}

	ctx := r.Context()

	assetLoan := &store.AssetLoan{
		AssetName:           payload.AssetName,
		AssetID:             payload.AssetID,
		UserID:              payload.UserID,
		CheckoutDate:        payload.CheckoutDate,
		ExpectedCheckinDate: payload.ExpectedCheckinDate,
		Notes:               payload.Notes,
	}

	if err := app.store.AssetLoan.Create(ctx, assetLoan); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.store.Asset.UpdateStatus(ctx, asset.ID, store.AssetAssigned); err != nil {
		app.internalServerError(w, r, err)
	}

	app.jsonResponse(w, http.StatusCreated, map[string]string{
		"status": "asset assigned successfully",
	})
}

func (app *application) assetContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "assetID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		asset, err := app.store.Asset.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, assetCtx, asset)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAssetFromCtx(r *http.Request) *store.Asset {
	asset, _ := r.Context().Value(assetCtx).(*store.Asset)
	return asset
}

func (app *application) updateAsset(ctx context.Context, asset *store.Asset) error {
	if err := app.store.Asset.Update(ctx, asset); err != nil {
		return err
	}

	return nil
}
