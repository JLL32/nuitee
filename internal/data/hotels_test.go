package data

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHotelModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	hotel := &Hotel{
		HotelID:     123,
		MainImageTh: "image.jpg",
		HotelName:   "Test Hotel",
		Phone:       "123-456-7890",
		Email:       "test@hotel.com",
		Address: Address{
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
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`INSERT INTO hotels`).
		WithArgs(
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
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	err = hotelModel.Insert(hotel)

	if err != nil {
		t.Errorf("error was not expected while inserting hotel: %s", err)
	}

	if hotel.CreatedAt != createdAt {
		t.Errorf("expected CreatedAt to be %v, got %v", createdAt, hotel.CreatedAt)
	}

	if hotel.UpdatedAt != updatedAt {
		t.Errorf("expected UpdatedAt to be %v, got %v", updatedAt, hotel.UpdatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_Insert_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	hotel := &Hotel{
		HotelID:   123,
		HotelName: "Test Hotel",
	}

	mock.ExpectQuery(`INSERT INTO hotels`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	err = hotelModel.Insert(hotel)

	if err == nil {
		t.Error("expected error, but got none")
	}

	if err != sql.ErrConnDone {
		t.Errorf("expected error to be %v, got %v", sql.ErrConnDone, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	expectedHotel := &Hotel{
		HotelID:     123,
		MainImageTh: "image.jpg",
		HotelName:   "Test Hotel",
		Phone:       "123-456-7890",
		Email:       "test@hotel.com",
		Address: Address{
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

	hotel, err := hotelModel.Get(123)

	if err != nil {
		t.Errorf("error was not expected while getting hotel: %s", err)
	}

	if !reflect.DeepEqual(hotel, expectedHotel) {
		t.Errorf("expected hotel %+v, got %+v", expectedHotel, hotel)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	mock.ExpectQuery(`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE hotel_id = \$1`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	hotel, err := hotelModel.Get(999)

	if hotel != nil {
		t.Error("expected hotel to be nil")
	}

	if err != ErrRecordNotFound {
		t.Errorf("expected error to be %v, got %v", ErrRecordNotFound, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_Get_InvalidID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	testCases := []int64{0, -1, -999}

	for _, id := range testCases {
		t.Run("invalid_id", func(t *testing.T) {
			hotel, err := hotelModel.Get(id)

			if hotel != nil {
				t.Error("expected hotel to be nil")
			}

			if err != ErrRecordNotFound {
				t.Errorf("expected error to be %v, got %v", ErrRecordNotFound, err)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "hotel_id",
		SortSafelist: []string{"hotel_id", "name"},
	}

	expectedHotels := []*Hotel{
		{
			HotelID:     123,
			MainImageTh: "image1.jpg",
			HotelName:   "Test Hotel 1",
			Phone:       "123-456-7890",
			Email:       "test1@hotel.com",
			Address: Address{
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
			Address: Address{
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
		WithArgs("test search", 20, 0).
		WillReturnRows(rows)

	hotels, metadata, err := hotelModel.GetAll("test search", filters)

	if err != nil {
		t.Errorf("error was not expected while getting all hotels: %s", err)
	}

	if len(hotels) != 2 {
		t.Errorf("expected 2 hotels, got %d", len(hotels))
	}

	if metadata.TotalRecords != 2 {
		t.Errorf("expected TotalRecords to be 2, got %d", metadata.TotalRecords)
	}

	if metadata.CurrentPage != 1 {
		t.Errorf("expected CurrentPage to be 1, got %d", metadata.CurrentPage)
	}

	if metadata.PageSize != 20 {
		t.Errorf("expected PageSize to be 20, got %d", metadata.PageSize)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_GetAll_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "hotel_id",
		SortSafelist: []string{"hotel_id", "name"},
	}

	rows := sqlmock.NewRows([]string{
		"count", "hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
		"city", "state", "country", "postal_code", "stars", "rating",
		"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
	})

	mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
		WithArgs("", 20, 0).
		WillReturnRows(rows)

	hotels, metadata, err := hotelModel.GetAll("", filters)

	if err != nil {
		t.Errorf("error was not expected while getting all hotels: %s", err)
	}

	if len(hotels) != 0 {
		t.Errorf("expected 0 hotels, got %d", len(hotels))
	}

	if metadata.TotalRecords != 0 {
		t.Errorf("expected TotalRecords to be 0, got %d", metadata.TotalRecords)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_GetAll_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "hotel_id",
		SortSafelist: []string{"hotel_id", "name"},
	}

	mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE fts @@ plainto_tsquery\('simple', \$1\) OR \$1 = '' ORDER BY hotel_id ASC, hotel_id ASC LIMIT \$2 OFFSET \$3`).
		WithArgs("", 20, 0).
		WillReturnError(sql.ErrConnDone)

	hotels, metadata, err := hotelModel.GetAll("", filters)

	if err == nil {
		t.Error("expected error, but got none")
	}

	if err != sql.ErrConnDone {
		t.Errorf("expected error to be %v, got %v", sql.ErrConnDone, err)
	}

	if hotels != nil {
		t.Error("expected hotels to be nil")
	}

	if metadata.TotalRecords != 0 {
		t.Errorf("expected empty metadata, got %+v", metadata)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Benchmark tests
func TestHotelModel_Upsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	hotel := &Hotel{
		HotelID:     123,
		MainImageTh: "image.jpg",
		HotelName:   "Test Hotel",
		Phone:       "123-456-7890",
		Email:       "test@hotel.com",
		Address: Address{
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
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`INSERT INTO hotels \(\s*hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description\s*\)\s*VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13, \$14, \$15, \$16\)\s*ON CONFLICT \(hotel_id\) DO UPDATE SET.*RETURNING created_at, updated_at`).
		WithArgs(
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
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	err = hotelModel.Upsert(hotel)

	if err != nil {
		t.Errorf("error was not expected while upserting hotel: %s", err)
	}

	if hotel.CreatedAt != createdAt {
		t.Errorf("expected CreatedAt to be %v, got %v", createdAt, hotel.CreatedAt)
	}

	if hotel.UpdatedAt != updatedAt {
		t.Errorf("expected UpdatedAt to be %v, got %v", updatedAt, hotel.UpdatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHotelModel_Upsert_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	hotel := &Hotel{
		HotelID:   123,
		HotelName: "Test Hotel",
	}

	mock.ExpectQuery(`INSERT INTO hotels \(\s*hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description\s*\)\s*VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13, \$14, \$15, \$16\)\s*ON CONFLICT \(hotel_id\) DO UPDATE SET.*RETURNING created_at, updated_at`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	err = hotelModel.Upsert(hotel)

	if err == nil {
		t.Error("expected error, but got none")
	}

	if err != sql.ErrConnDone {
		t.Errorf("expected error to be %v, got %v", sql.ErrConnDone, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func BenchmarkHotelModel_Get(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hotelModel := HotelModel{DB: db}

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(`SELECT hotel_id, main_image_th, hotel_name, phone, email, address, city, state, country, postal_code, stars, rating, review_count, child_allowed, pets_allowed, description, created_at, updated_at FROM hotels WHERE hotel_id = \$1`).
			WithArgs(int64(123)).
			WillReturnRows(sqlmock.NewRows([]string{
				"hotel_id", "main_image_th", "hotel_name", "phone", "email", "address",
				"city", "state", "country", "postal_code", "stars", "rating",
				"review_count", "child_allowed", "pets_allowed", "description", "created_at", "updated_at",
			}).AddRow(
				123, "image.jpg", "Test Hotel", "123-456-7890", "test@hotel.com", "123 Main St",
				"Test City", "Test State", "Test Country", "12345", 5, 4.5,
				100, true, false, "A wonderful test hotel", time.Now(), time.Now(),
			))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := hotelModel.Get(123)
		if err != nil {
			b.Errorf("unexpected error: %s", err)
		}
	}
}
