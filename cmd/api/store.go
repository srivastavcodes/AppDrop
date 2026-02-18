package main

import (
	"appdrop/internal/data"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (b *backend) listStoresHandler(w http.ResponseWriter, r *http.Request) {
	stores, err := b.models.Stores.GetAll()
	if err != nil {
		b.serverErrorResponse(w, r, err)
		return
	}
	if err = b.writeJson(w, http.StatusOK, envelope{"stores": stores}, nil); err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

func (b *backend) showStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r, "store_id")
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	store, err := b.models.Stores.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			b.notFoundResponse(w, r)
		} else {
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	if err = b.writeJson(w, http.StatusOK, envelope{"store": store}, nil); err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

func (b *backend) createStoreHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := b.readJson(w, r, &input); err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	if strings.TrimSpace(input.Name) == "" {
		b.validationErrorResponse(w, r, "store name is required")
		return
	}
	if strings.TrimSpace(input.Slug) == "" {
		b.validationErrorResponse(w, r, "store slug is required")
		return
	}

	store := &data.Store{
		Id:   uuid.New(),
		Name: input.Name,
		Slug: input.Slug,
	}
	if err := b.models.Stores.Insert(store); err != nil {
		if strings.Contains(err.Error(), "slug already exists") {
			b.conflictResponse(w, r, err.Error())
		} else {
			b.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/stores/%s", store.Id))

	if err := b.writeJson(w, http.StatusCreated, envelope{"store": store}, headers); err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

func (b *backend) updateStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r, "store_id")
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	store, err := b.models.Stores.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			b.notFoundResponse(w, r)
		} else {
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Name *string `json:"name"`
		Slug *string `json:"slug"`
	}
	if err = b.readJson(w, r, &input); err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			b.validationErrorResponse(w, r, "store name cannot be empty")
			return
		}
		store.Name = *input.Name
	}
	if input.Slug != nil {
		if strings.TrimSpace(*input.Slug) == "" {
			b.validationErrorResponse(w, r, "store slug cannot be empty")
			return
		}
		store.Slug = *input.Slug
	}

	if err = b.models.Stores.Update(store); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		case strings.Contains(err.Error(), "slug already exists"):
			b.conflictResponse(w, r, err.Error())
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	if err = b.writeJson(w, http.StatusOK, envelope{"store": store}, nil); err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

func (b *backend) deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r, "store_id")
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	if err = b.models.Stores.Delete(id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"message": "store successfully deleted"}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}
