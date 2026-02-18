package main

import (
	"appdrop/internal/data"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// todo: put validation logic separate later

// createPageHandler handles POST /pages
func (b *backend) createPageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AppId  uuid.UUID `json:"app_id"`
		Name   string    `json:"name"`
		Route  string    `json:"route"`
		IsHome bool      `json:"is_home"`
	}
	err := b.readJson(w, r, &input)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	// Validate input
	if strings.TrimSpace(input.Name) == "" {
		b.validationErrorResponse(w, r, "page name is required and cannot be empty")
		return
	}
	if strings.TrimSpace(input.Route) == "" {
		b.validationErrorResponse(w, r, "page route is required and cannot be empty")
		return
	}
	// For this assignment, we'll use a default app_id if not provided
	// In production, this would come from the authentication context
	if input.AppId == uuid.Nil {
		input.AppId = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}
	page := &data.Page{
		Id: uuid.New(), StoreId: input.AppId,
		IsHome: input.IsHome,
		Name:   input.Name, Route: input.Route,
	}
	err = b.models.Pages.Insert(page)
	if err != nil {
		if strings.Contains(err.Error(), "page route already exists") {
			b.conflictResponse(w, r, "page route already exists")
		} else {
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/pages/%v", page.Id))

	err = b.writeJson(w, http.StatusCreated, envelope{"page": page}, headers)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// listPagesHandler handles GET /pages
func (b *backend) listPagesHandler(w http.ResponseWriter, r *http.Request) {
	pages, err := b.models.Pages.GetAll()
	if err != nil {
		b.serverErrorResponse(w, r, err)
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"pages": pages}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// listPagesHandler handles GET /pages
func (b *backend) listPagesForAppHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	pages, err := b.models.Pages.GetAllForApp(id)
	if err != nil {
		b.serverErrorResponse(w, r, err)
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"pages": pages}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// showPageHandler handles GET /pages/:id
func (b *backend) showPageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	page, err := b.models.Pages.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"page": page}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// updatePageHandler handles PUT /pages/:id
func (b *backend) updatePageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	page, err := b.models.Pages.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Name   *string `json:"name"`
		Route  *string `json:"route"`
		IsHome *bool   `json:"is_home"`
	}
	if err = b.readJson(w, r, &input); err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	// Update fields if provided
	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			b.validationErrorResponse(w, r, "page name cannot be empty")
			return
		}
		page.Name = *input.Name
	}
	if input.Route != nil {
		if strings.TrimSpace(*input.Route) == "" {
			b.validationErrorResponse(w, r, "page route cannot be empty")
			return
		}
		page.Route = *input.Route
	}
	if input.IsHome != nil {
		page.IsHome = *input.IsHome
	}

	err = b.models.Pages.Update(page)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		case strings.Contains(err.Error(), "page route already exists"):
			b.conflictResponse(w, r, "page route already exists for this b")
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"page": page}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// deletePageHandler handles DELETE /pages/:id
func (b *backend) deletePageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	err = b.models.Pages.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		case strings.Contains(err.Error(), "cannot delete home page"):
			b.conflictResponse(w, r, "cannot delete home page")
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"message": "page successfully deleted"}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}
