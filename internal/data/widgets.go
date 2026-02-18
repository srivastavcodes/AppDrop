package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Widget struct {
	Id        uuid.UUID      `json:"id"`
	PageId    uuid.UUID      `json:"page_id"`
	Type      string         `json:"type"`
	Position  int            `json:"position"`
	Config    map[string]any `json:"config,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type WidgetModel struct {
	Db *sql.DB
}

// GetForPage returns all widgets for a specific page, ordered by position
func (m *WidgetModel) GetForPage(pageID uuid.UUID) ([]*Widget, error) {
	query := `SELECT id, page_id, type, position, config, created_at, updated_at FROM widgets
		    WHERE page_id = $1 ORDER BY position`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Db.QueryContext(ctx, query, pageID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows:", "err", err)
		}
	}()
	var widgets []*Widget

	for rows.Next() {
		var configJSON []byte
		var widget Widget

		err := rows.Scan(
			&widget.Id, &widget.PageId,
			&widget.Type, &widget.Position,
			&configJSON,
			&widget.CreatedAt, &widget.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if configJSON != nil {
			err = json.Unmarshal(configJSON, &widget.Config)
			if err != nil {
				return nil, err
			}
		}
		widgets = append(widgets, &widget)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return widgets, nil
}

// Insert creates a new widget
func (m *WidgetModel) Insert(widget *Widget) error {
	// Get the next position for this page
	posQuery := `SELECT COALESCE(MAX(position), -1) + 1 FROM widgets WHERE page_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.Db.QueryRowContext(ctx, posQuery, widget.PageId).Scan(&widget.Position)
	if err != nil {
		return err
	}
	// Marshal config to JSON
	var configJSON []byte
	if widget.Config != nil {
		configJSON, err = json.Marshal(widget.Config)
		if err != nil {
			return err
		}
	}
	query := `INSERT INTO widgets (id, page_id, type, position, config) VALUES ($1, $2, $3, $4, $5) 
		    RETURNING created_at, updated_at`

	args := []any{
		widget.Id, widget.PageId,
		widget.Type,
		widget.Position, configJSON,
	}
	err = m.Db.QueryRowContext(ctx, query, args...).Scan(&widget.CreatedAt, &widget.UpdatedAt)
	return err
}

// Get returns a single widget by ID.
func (m *WidgetModel) Get(id uuid.UUID) (*Widget, error) {
	query := `SELECT id, page_id, type, position, config, created_at, updated_at FROM widgets 
		    WHERE id = $1`
	var widget Widget

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var configJSON []byte

	err := m.Db.QueryRowContext(ctx, query, id).Scan(
		&widget.Id, &widget.PageId,
		&widget.Type, &widget.Position,
		&configJSON,
		&widget.CreatedAt, &widget.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	if configJSON != nil {
		if err = json.Unmarshal(configJSON, &widget.Config); err != nil {
			return nil, err
		}
	}
	return &widget, nil
}

// Update modifies an existing widget.
func (m *WidgetModel) Update(widget *Widget) error {
	var configJSON []byte
	var err error

	if widget.Config != nil {
		configJSON, err = json.Marshal(widget.Config)
		if err != nil {
			return err
		}
	}
	query := `UPDATE widgets SET type = $1, position = $2, config = $3, updated_at = $4 WHERE id = $5 
		    RETURNING updated_at`

	args := []any{widget.Type, widget.Position, configJSON, widget.UpdatedAt, widget.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = m.Db.QueryRowContext(ctx, query, args...).Scan(&widget.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

// Delete removes a widget.
func (m *WidgetModel) Delete(id uuid.UUID) error {
	query := `DELETE FROM widgets WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.Db.ExecContext(ctx, query, id)
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

// Reorder updates the positions of multiple widgets.
func (m *WidgetModel) Reorder(pageID uuid.UUID, widgetIDs []uuid.UUID) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Verify all widgets belong to this page
	verifyQuery := `SELECT COUNT(*) FROM widgets WHERE page_id = $1 AND id = ANY($2)`

	var count int
	err = tx.QueryRowContext(ctx, verifyQuery, pageID, pq.Array(widgetIDs)).Scan(&count)
	if err != nil {
		return err
	}
	if count != len(widgetIDs) {
		return errors.New("some widgets do not belong to this page")
	}
	updateQuery := `UPDATE widgets SET position = $1 WHERE id = $2 AND page_id = $3`

	for i, widgetID := range widgetIDs {
		_, err = tx.ExecContext(ctx, updateQuery, i, widgetID, pageID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
