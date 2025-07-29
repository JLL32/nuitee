package data

import (
	"context"
	"database/sql"
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

/*
CREATE TABLE hotels (

	hotel_id INTEGER PRIMARY KEY,
	main_image_th TEXT,
	hotel_name TEXT NOT NULL,
	phone TEXT,
	email TEXT,
	address TEXT,
	city TEXT,
	state TEXT,
	country TEXT,
	postal_code TEXT,
	stars INTEGER,
	rating DECIMAL(3,2),
	review_count INTEGER,
	child_allowed BOOLEAN,
	pets_allowed BOOLEAN,
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);
*/
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
