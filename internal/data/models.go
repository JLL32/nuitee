package data

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Hotels HotelModel
	Reviews ReviewModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Hotels: HotelModel{DB: db},
		Reviews: ReviewModel{DB: db},
	}
}
