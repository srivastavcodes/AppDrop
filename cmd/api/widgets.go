package main

import (
	"appdrop/internal/data"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Allowed widget types as per requirements
var allowedWidgetTypes = map[string]bool{
	"banner":       true,
	"product_grid": true,
	"text":         true,
	"image":        true,
	"spacer":       true,
}

// createWidgetHandler handles POST /pages/:id/widgets
func (b *backend) createWidgetHandler(w http.ResponseWriter, r *http.Request) {
	pageID, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	// Verify page exists
	_, err = b.models.Pages.Get(pageID)
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
		Type   string         `json:"type"`
		Config map[string]any `json:"config"`
	}
	err = b.readJson(w, r, &input)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	// Validate widget type
	if !allowedWidgetTypes[input.Type] {
		b.validationErrorResponse(w, r, "widget type must be one of: banner, product_grid, text, image, spacer")
		return
	}
	widget := &data.Widget{
		Id:     uuid.New(),
		PageId: pageID,
		Type:   input.Type,
		Config: input.Config,
	}
	err = b.models.Widgets.Insert(widget)
	if err != nil {
		b.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("widgets/%s", widget.Id.String()))

	err = b.writeJson(w, http.StatusCreated, envelope{"widget": widget}, headers)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// updateWidgetHandler handles PUT /widgets/:id
func (b *backend) updateWidgetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	var input struct {
		Type   *string         `json:"type"`
		Config *map[string]any `json:"config"`
	}
	err = b.readJson(w, r, &input)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	widget, err := b.models.Widgets.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		default:
			b.serverErrorResponse(w, r, err)
		}
	}
	if input.Type != nil {
		if !allowedWidgetTypes[*input.Type] {
			b.validationErrorResponse(
				w, r,
				"widget type must be one of: banner, product_grid, text, image, spacer",
			)
			return
		}
		widget.Type = *input.Type
	}
	if input.Config != nil {
		widget.Config = *input.Config
	}
	widget.UpdatedAt = time.Now()

	err = b.models.Widgets.Update(widget)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"widget": widget}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// deleteWidgetHandler handles DELETE /widgets/:id
func (b *backend) deleteWidgetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	err = b.models.Widgets.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			b.notFoundResponse(w, r)
		default:
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"message": "widget successfully deleted"}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}

// reorderWidgetsHandler handles POST /pages/:id/widgets/reorder
func (b *backend) reorderWidgetsHandler(w http.ResponseWriter, r *http.Request) {
	pageID, err := b.readIdParam(r)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	// Verify page exists
	_, err = b.models.Pages.Get(pageID)
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
		WidgetIds []uuid.UUID `json:"widget_ids"`
	}
	err = b.readJson(w, r, &input)
	if err != nil {
		b.badRequestResponse(w, r, err)
		return
	}
	if len(input.WidgetIds) == 0 {
		b.validationErrorResponse(w, r, "widget_ids array cannot be empty")
		return
	}

	err = b.models.Widgets.Reorder(pageID, input.WidgetIds)
	if err != nil {
		if err.Error() == "some widgets do not belong to this page" {
			b.validationErrorResponse(w, r, err.Error())
		} else {
			b.serverErrorResponse(w, r, err)
		}
		return
	}
	err = b.writeJson(w, http.StatusOK, envelope{"message": "widgets successfully reordered"}, nil)
	if err != nil {
		b.serverErrorResponse(w, r, err)
	}
}
