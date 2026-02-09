package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Page struct {
	Id        uuid.UUID `json:"id"`
	AppId     uuid.UUID `json:"app_id"`
	Name      string    `json:"name"`
	Route     string    `json:"route"`
	IsHome    bool      `json:"is_home"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Widgets   []*Widget `json:"widgets,omitempty"`
}

type PageModel struct {
	Db *sql.DB
}

// Insert creates a new page and returns created at and updated at from db.
func (pm *PageModel) Insert(page *Page) error {
	query := `INSERT INTO pages (id, app_id, name, route, is_home) VALUES ($1, $2, $3, $4, $5) 
		    RETURNING created_at, updated_at`

	args := []any{page.Id, page.AppId, page.Name, page.Route, page.IsHome}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If this page is set as home, unset all other home pages for this app
	if page.IsHome {
		err := pm.changeIsHomePage(ctx, page.AppId, nil)
		if err != nil {
			return err
		}
	}
	err := pm.Db.QueryRowContext(ctx, query, args...).Scan(&page.CreatedAt, &page.UpdatedAt)
	if err != nil {
		switch {
		case err.(*pq.Error).Code.Name() == "unique_violation":
			return errors.New("page route already exists for this app")
		default:
			return err
		}
	}
	return nil
}

// GetAll returns all pages for a specific app.
func (pm *PageModel) GetAll() ([]*Page, error) {
	query := `SELECT id, app_id, name, route, is_home, created_at, updated_at FROM pages 
		    ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := pm.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows:", "err", err)
		}
	}()
	var pages []*Page

	for rows.Next() {
		var page Page

		err := rows.Scan(
			&page.Id, &page.AppId,
			&page.Name, &page.Route,
			&page.IsHome,
			&page.CreatedAt, &page.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pages = append(pages, &page)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}

// GetAllForApp returns all pages for a specific app.
func (pm *PageModel) GetAllForApp(appID uuid.UUID) ([]*Page, error) {
	if appID == uuid.Nil {
		return nil, errors.New("appID is required")
	}
	query := `SELECT id, app_id, name, route, is_home, created_at, updated_at FROM pages 
		    WHERE app_id = $1 ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := pm.Db.QueryContext(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows:", "err", err)
		}
	}()
	var pages []*Page

	for rows.Next() {
		var page Page

		err := rows.Scan(
			&page.Id, &page.AppId,
			&page.Name, &page.Route,
			&page.IsHome,
			&page.CreatedAt, &page.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pages = append(pages, &page)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}

// Get returns a single page with its widgets.
func (pm *PageModel) Get(id uuid.UUID) (*Page, error) {
	query := `SELECT id, app_id, name, route, is_home, created_at, updated_at
		    FROM pages WHERE id = $1`

	var page Page

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pm.Db.QueryRowContext(ctx, query, id).Scan(
		&page.Id, &page.AppId,
		&page.Name, &page.Route,
		&page.IsHome,
		&page.CreatedAt, &page.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	wm := WidgetModel{Db: pm.Db}

	w, err := wm.GetForPage(page.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting page widgets: %w", err)
	}
	page.Widgets = w

	return &page, nil
}

// Update modifies an existing page properties.
func (pm *PageModel) Update(page *Page) error {
	query := `UPDATE pages SET name = $1, route = $2, is_home = $3 WHERE id = $4 
		    RETURNING updated_at`

	args := []any{page.Name, page.Route, page.IsHome, page.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If setting this as home, unset all others for this app first
	if page.IsHome {
		err := pm.changeIsHomePage(ctx, page.AppId, &page.Id)
		if err != nil {
			return err
		}
	}
	err := pm.Db.QueryRowContext(ctx, query, args...).Scan(&page.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case err.(*pq.Error).Code.Name() == "unique_violation":
			return errors.New("page route already exists for this app")
		default:
			return err
		}
	}
	return nil
}

// Delete removes a page and all its widgets (CASCADE)
func (pm *PageModel) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}
	checkQuery := `SELECT is_home FROM pages WHERE id = $1`
	var isHome bool

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pm.Db.QueryRowContext(ctx, checkQuery, id).Scan(&isHome)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	if isHome {
		return errors.New("cannot delete home page")
	}
	query := `DELETE FROM pages WHERE id = $1`

	result, err := pm.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// changeIsHomePage sets is_home to false for all pages except the specified one
func (pm *PageModel) changeIsHomePage(ctx context.Context, appID uuid.UUID, exceptID *uuid.UUID) error {
	var query string
	var args []any

	if exceptID == nil {
		// Unset all home pages for this app
		query = `UPDATE pages SET is_home = FALSE WHERE app_id = $1 AND is_home = TRUE`
		args = []any{appID}
	} else {
		// Unset all except the specified page
		query = `UPDATE pages SET is_home = FALSE WHERE app_id = $1 AND is_home = TRUE AND id != $2`
		args = []any{appID, *exceptID}
	}
	_, err := pm.Db.ExecContext(ctx, query, args...)
	return err
}
