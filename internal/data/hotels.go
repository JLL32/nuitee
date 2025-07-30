package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Address struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

type Hotel struct {
	HotelID      int       `json:"hotel_id"`
	MainImageTh  string    `json:"main_image_th"`
	HotelName    string    `json:"hotel_name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Address      Address   `json:"address"`
	Stars        int       `json:"stars"`
	Rating       float64   `json:"rating"`
	ReviewCount  int       `json:"review_count"`
	ChildAllowed bool      `json:"child_allowed"`
	PetsAllowed  bool      `json:"pets_allowed"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type HotelModel struct {
	DB *sql.DB
}

func (h HotelModel) Insert(hotel *Hotel) error {
	query :=
		`INSERT INTO hotels (
			hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at, updated_at`

	args := []any{
		hotel.HotelID,
		hotel.MainImageTh,
		hotel.HotelName,
		hotel.Phone,
		hotel.Email,
		hotel.Address.Address,
		hotel.Address.City,
		hotel.Address.State,
		hotel.Address.Country,
		hotel.Address.PostalCode,
		hotel.Stars,
		hotel.Rating,
		hotel.ReviewCount,
		hotel.ChildAllowed,
		hotel.PetsAllowed,
		hotel.Description,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return h.DB.QueryRowContext(ctx, query, args...).Scan(&hotel.CreatedAt, &hotel.UpdatedAt)
}

func (h HotelModel) Get(id int64) (*Hotel, error) {
	if id <= 0 {
		return nil, ErrRecordNotFound
	}

	query :=
		`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at
		FROM hotels
		WHERE hotel_id = $1`

	var hotel Hotel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := h.DB.QueryRowContext(ctx, query, id).Scan(
		&hotel.HotelID,
		&hotel.MainImageTh,
		&hotel.HotelName,
		&hotel.Phone,
		&hotel.Email,
		&hotel.Address.Address,
		&hotel.Address.City,
		&hotel.Address.State,
		&hotel.Address.Country,
		&hotel.Address.PostalCode,
		&hotel.Stars,
		&hotel.Rating,
		&hotel.ReviewCount,
		&hotel.ChildAllowed,
		&hotel.PetsAllowed,
		&hotel.Description,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &hotel, nil
}
