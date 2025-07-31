package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

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

func (r ReviewModel) Get(hotelID int64, id int64) (*Review, error) {
	if id <= 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at
		FROM reviews
		WHERE id = $1 AND hotel_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var review Review
	err := r.DB.QueryRowContext(ctx, query, id, hotelID).Scan(
		&review.ID,
		&review.HotelID,
		&review.AverageScore,
		&review.Country,
		&review.Type,
		&review.Name,
		&review.Date,
		&review.Headline,
		&review.Language,
		&review.Pros,
		&review.Cons,
		&review.Source,
		&review.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &review, nil
}

func (r ReviewModel) GetAll(search string, filters Filters) ([]*Review, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at
		FROM reviews
		WHERE fts @@ plainto_tsquery('simple', $1) OR $1 = ''
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, search, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRows := 0
	var reviews []*Review

	for rows.Next() {
		var review Review

		err := rows.Scan(
			&totalRows,
			&review.ID,
			&review.HotelID,
			&review.AverageScore,
			&review.Country,
			&review.Type,
			&review.Name,
			&review.Date,
			&review.Headline,
			&review.Language,
			&review.Pros,
			&review.Cons,
			&review.Source,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRows, filters.Page, filters.PageSize)

	return reviews, metadata, nil
}
