package data

import (
	"context"
	"database/sql"
	"time"
)

/*
*
CREATE TABLE reviews (

	id SERIAL PRIMARY KEY,
	hotel_id INTEGER NOT NULL,
	average_score INTEGER,
	country CHAR(2),
	type TEXT,
	name TEXT,
	date TIMESTAMP,
	headline TEXT,
	language CHAR(2),
	pros TEXT,
	cons TEXT,
	source TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (hotel_id) REFERENCES hotels(hotel_id) ON DELETE CASCADE

);
*/
type Review struct {
	ID           int       `json:"id"`
	HotelID      int       `json:"hotel_id"`
	AverageScore int       `json:"average_score"`
	Country      string    `json:"country"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Date         string    `json:"date"`
	Headline     string    `json:"headline"`
	Language     string    `json:"language"`
	Pros         string    `json:"pros"`
	Cons         string    `json:"cons"`
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"created_at"`
}

type ReviewModel struct {
	DB *sql.DB
}

func NewReviewModel(db *sql.DB) *ReviewModel {
	return &ReviewModel{DB: db}
}

func (r ReviewModel) Insert(hotelID int, review *Review) error {
	query := `
		INSERT INTO reviews (hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, hotel_id, created_at
	`

	args := []any{
		hotelID,
		review.AverageScore,
		review.Country,
		review.Type,
		review.Name,
		review.Date,
		review.Headline,
		review.Language,
		review.Pros,
		review.Cons,
		review.Source,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.HotelID, &review.CreatedAt)
}
