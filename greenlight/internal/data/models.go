package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	// Movie最好与实际表名对应
	Movie MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movie: MovieModel{DB: db},
	}
}
