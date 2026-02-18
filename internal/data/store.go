package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Store struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StoreModel struct {
	Db *sql.DB
}

func (m *StoreModel) Insert(store *Store) error {
	query := `INSERT INTO stores (id, name, slug) VALUES ($1, $2, $3) RETURNING created_at, updated_at`

	args := []any{store.Id, store.Name, store.Slug}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.Db.QueryRowContext(ctx, query, args...).Scan(
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return errors.New("store slug already exists")
		}
		return err
	}
	return nil
}

func (m *StoreModel) Get(id uuid.UUID) (*Store, error) {
	// todo: get pages here as well?

	query := `SELECT id, name, slug, created_at, updated_at FROM stores WHERE id = $1`

	var store Store

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.Db.QueryRowContext(ctx, query, id).Scan(
		&store.Id, &store.Name,
		&store.Slug,
		&store.CreatedAt, &store.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &store, nil
}

func (m *StoreModel) GetAll() ([]*Store, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM stores ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "err", err)
		}
	}()
	var stores []*Store

	for rows.Next() {
		var store Store

		err := rows.Scan(
			&store.Id, &store.Name,
			&store.Slug,
			&store.CreatedAt, &store.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return stores, nil
}

func (m *StoreModel) Update(store *Store) error {
	query := `UPDATE stores SET name = $1, slug = $2, updated_at = NOW() WHERE id = $3 
		    RETURNING updated_at`

	args := []any{store.Name, store.Slug, store.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.Db.QueryRowContext(ctx, query, args...).Scan(&store.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case func() bool {
			var pqErr *pq.Error
			ok := errors.As(err, &pqErr)
			return ok && pqErr.Code.Name() == "unique_violation"
		}():
			return errors.New("store slug already exists")
		default:
			return err
		}
	}
	return nil
}

func (m *StoreModel) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.Db.ExecContext(ctx, `DELETE FROM stores WHERE id = $1`, id)
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
