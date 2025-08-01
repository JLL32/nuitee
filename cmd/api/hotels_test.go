package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JLL32/nuitee/internal/data"
	"github.com/julienschmidt/httprouter"
)



func TestGetHotelHandler(t *testing.T) {
	app, mock, cleanup := newTestApplication(t)
	defer cleanup()

	// Test data
	expectedHotel := &data.Hotel{
		HotelID:     123,
		MainImageTh: "image.jpg",
		HotelName:   "Test Hotel",
		Phone:       "123-456-7890",
		Email:       "test@hotel.com",
		Address: data.Address{
			Address:    "123 Main St",
			City:       "Test City",
			State:      "Test State",
			Country:    "Test Country",
			PostalCode: "12345",
		},
		Stars:        5,
		Rating:       4.5,
		ReviewCount:  100,
		ChildAllowed: true,
		PetsAllowed:  false,
		Description:  "A wonderful test hotel",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name           string
		hotelID        string
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "valid hotel ID",
			hotelID: "123",
			setupMock: func() {
				mock.ExpectQuery(`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE hotel_id = \$1`).
					WithArgs(int64(123)).
					WillReturnRows(sqlmock.NewRows([]string{
						"hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
						"city", "state", "country", "postal_code", "stars", "rating",
						"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
					}).AddRow(
						expectedHotel.HotelID, expectedHotel.MainImageTh, expectedHotel.HotelName,
						expectedHotel.Phone, expectedHotel.Email, expectedHotel.Address.Address,
						expectedHotel.Address.City, expectedHotel.Address.State, expectedHotel.Address.Country,
						expectedHotel.Address.PostalCode, expectedHotel.Stars, expectedHotel.Rating,
						expectedHotel.ReviewCount, expectedHotel.ChildAllowed, expectedHotel.PetsAllowed,
						expectedHotel.Description, expectedHotel.CreatedAt, expectedHotel.UpdatedAt,
					))
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response struct {
					Hotel data.Hotel `json:"hotel"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response.Hotel.HotelID != expectedHotel.HotelID {
					t.Errorf("expected hotel ID %d, got %d", expectedHotel.HotelID, response.Hotel.HotelID)
				}

				if response.Hotel.HotelName != expectedHotel.HotelName {
					t.Errorf("expected hotel name %s, got %s", expectedHotel.HotelName, response.Hotel.HotelName)
				}
			},
		},
		{
			name:    "hotel not found",
			hotelID: "999",
			setupMock: func() {
				mock.ExpectQuery(`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE hotel_id = \$1`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
		{
			name:    "invalid hotel ID",
			hotelID: "invalid",
			setupMock: func() {
				// No mock setup needed as validation should fail before DB call
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
		{
			name:    "database error",
			hotelID: "123",
			setupMock: func() {
				mock.ExpectQuery(`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE hotel_id = \$1`).
					WithArgs(int64(123)).
					WillReturnError(sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req, err := http.NewRequest("GET", fmt.Sprintf("/v1/hotels/%s", tt.hotelID), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			// Create router and add route parameter
			router := httprouter.New()
			router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID", app.getHotelHandler)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			tt.checkResponse(t, rr)

			// Verify all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListHotelsHandler(t *testing.T) {
	app, mock, cleanup := newTestApplication(t)
	defer cleanup()

	// Test data
	expectedHotels := []*data.Hotel{
		{
			HotelID:     123,
			MainImageTh: "image1.jpg",
			HotelName:   "Test Hotel 1",
			Phone:       "123-456-7890",
			Email:       "test1@hotel.com",
			Address: data.Address{
				Address:    "123 Main St",
				City:       "Test City",
				State:      "Test State",
				Country:    "Test Country",
				PostalCode: "12345",
			},
			Stars:        5,
			Rating:       4.5,
			ReviewCount:  100,
			ChildAllowed: true,
			PetsAllowed:  false,
			Description:  "A wonderful test hotel 1",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			HotelID:     124,
			MainImageTh: "image2.jpg",
			HotelName:   "Test Hotel 2",
			Phone:       "123-456-7891",
			Email:       "test2@hotel.com",
			Address: data.Address{
				Address:    "124 Main St",
				City:       "Test City 2",
				State:      "Test State 2",
				Country:    "Test Country 2",
				PostalCode: "12346",
			},
			Stars:        4,
			Rating:       4.2,
			ReviewCount:  80,
			ChildAllowed: false,
			PetsAllowed:  true,
			Description:  "A wonderful test hotel 2",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "default parameters",
			queryParams: "",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"count", "hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
					"city", "state", "country", "postal_code", "stars", "rating",
					"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
				})

				for _, hotel := range expectedHotels {
					rows.AddRow(
						2, hotel.HotelID, hotel.MainImageTh, hotel.HotelName,
						hotel.Phone, hotel.Email, hotel.Address.Address,
						hotel.Address.City, hotel.Address.State, hotel.Address.Country,
						hotel.Address.PostalCode, hotel.Stars, hotel.Rating,
						hotel.ReviewCount, hotel.ChildAllowed, hotel.PetsAllowed,
						hotel.Description, hotel.CreatedAt, hotel.UpdatedAt,
					)
				}

				mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
					WithArgs("", 20, 0).
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response struct {
					Hotels   []data.Hotel   `json:"hotels"`
					Metadata data.Metadata `json:"metadata"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if len(response.Hotels) != 2 {
					t.Errorf("expected 2 hotels, got %d", len(response.Hotels))
				}

				if response.Metadata.TotalRecords != 2 {
					t.Errorf("expected TotalRecords to be 2, got %d", response.Metadata.TotalRecords)
				}
			},
		},
		{
			name:        "with search parameter",
			queryParams: "search=luxury",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"count", "hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
					"city", "state", "country", "postal_code", "stars", "rating",
					"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
				})

				// Add one hotel for search results
				hotel := expectedHotels[0]
				rows.AddRow(
					1, hotel.HotelID, hotel.MainImageTh, hotel.HotelName,
					hotel.Phone, hotel.Email, hotel.Address.Address,
					hotel.Address.City, hotel.Address.State, hotel.Address.Country,
					hotel.Address.PostalCode, hotel.Stars, hotel.Rating,
					hotel.ReviewCount, hotel.ChildAllowed, hotel.PetsAllowed,
					hotel.Description, hotel.CreatedAt, hotel.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
					WithArgs("luxury", 20, 0).
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response struct {
					Hotels   []data.Hotel   `json:"hotels"`
					Metadata data.Metadata `json:"metadata"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if len(response.Hotels) != 1 {
					t.Errorf("expected 1 hotel, got %d", len(response.Hotels))
				}
			},
		},
		{
			name:        "with pagination parameters",
			queryParams: "page=2&page_size=10",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"count", "hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
					"city", "state", "country", "postal_code", "stars", "rating",
					"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
					WithArgs("", 10, 10).
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response struct {
					Hotels   []data.Hotel   `json:"hotels"`
					Metadata data.Metadata `json:"metadata"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if len(response.Hotels) != 0 {
					t.Errorf("expected 0 hotels, got %d", len(response.Hotels))
				}
			},
		},
		{
			name:        "invalid page parameter",
			queryParams: "page=0",
			setupMock: func() {
				// No mock setup needed as validation should fail before DB call
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
		{
			name:        "invalid page_size parameter",
			queryParams: "page_size=101",
			setupMock: func() {
				// No mock setup needed as validation should fail before DB call
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
		{
			name:        "invalid sort parameter",
			queryParams: "sort=invalid_field",
			setupMock: func() {
				// No mock setup needed as validation should fail before DB call
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
		{
			name:        "database error",
			queryParams: "",
			setupMock: func() {
				mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
					WithArgs("", 20, 0).
					WillReturnError(sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["error"] == nil {
					t.Error("expected error field in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var reqURL string
			if tt.queryParams != "" {
				reqURL = fmt.Sprintf("/v1/hotels?%s", tt.queryParams)
			} else {
				reqURL = "/v1/hotels"
			}

			req, err := http.NewRequest("GET", reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			// Create router and add route
			router := httprouter.New()
			router.HandlerFunc(http.MethodGet, "/v1/hotels", app.listHotelsHandler)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			tt.checkResponse(t, rr)

			// Verify all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetHotelHandler(b *testing.B) {
	app, mock, cleanup := newTestApplicationForBenchmark(b)
	defer cleanup()

	// Setup mock for benchmark
	expectedHotel := &data.Hotel{
		HotelID:     123,
		MainImageTh: "image.jpg",
		HotelName:   "Test Hotel",
		Phone:       "123-456-7890",
		Email:       "test@hotel.com",
		Address: data.Address{
			Address:    "123 Main St",
			City:       "Test City",
			State:      "Test State",
			Country:    "Test Country",
			PostalCode: "12345",
		},
		Stars:        5,
		Rating:       4.5,
		ReviewCount:  100,
		ChildAllowed: true,
		PetsAllowed:  false,
		Description:  "A wonderful test hotel",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE hotel_id = \$1`).
			WithArgs(int64(123)).
			WillReturnRows(sqlmock.NewRows([]string{
				"hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
				"city", "state", "country", "postal_code", "stars", "rating",
				"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
			}).AddRow(
				expectedHotel.HotelID, expectedHotel.MainImageTh, expectedHotel.HotelName,
				expectedHotel.Phone, expectedHotel.Email, expectedHotel.Address.Address,
				expectedHotel.Address.City, expectedHotel.Address.State, expectedHotel.Address.Country,
				expectedHotel.Address.PostalCode, expectedHotel.Stars, expectedHotel.Rating,
				expectedHotel.ReviewCount, expectedHotel.ChildAllowed, expectedHotel.PetsAllowed,
				expectedHotel.Description, expectedHotel.CreatedAt, expectedHotel.UpdatedAt,
			))
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID", app.getHotelHandler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/v1/hotels/123", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}

func BenchmarkListHotelsHandler(b *testing.B) {
	app, mock, cleanup := newTestApplicationForBenchmark(b)
	defer cleanup()

	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{
			"count", "hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
			"city", "state", "country", "postal_code", "stars", "rating",
			"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
		}).AddRow(
			1, 123, "image.jpg", "Test Hotel", "123-456-7890", "test@hotel.com", "123 Main St",
			"Test City", "Test State", "Test Country", "12345", 5, 4.5,
			100, true, false, "A wonderful test hotel", time.Now(), time.Now(),
		)

		mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
			WithArgs("", 20, 0).
			WillReturnRows(rows)
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/hotels", app.listHotelsHandler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/v1/hotels", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}