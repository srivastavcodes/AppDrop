package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models groups the application’s data models behind a single dependency.
// It provides a convenient way to pass model access through handlers and
// services.
type Models struct {
	Pages   PageModel
	Widgets WidgetModel
}

// NewModels returns a new model with the fields initialized with the given db.
func NewModels(db *sql.DB) Models {
	return Models{
		Pages: PageModel{
			Db: db,
		},
		Widgets: WidgetModel{
			Db: db,
		},
	}
}
