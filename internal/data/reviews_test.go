package data

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestReviewModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	review := &Review{
		AverageScore: 8,
		Country:      "USA",
		Type:         "Business",
		Name:         "John Doe",
		Date:         "2024-01-15",
		Headline:     "Great stay!",
		Language:     "en",
		Pros:         "Clean rooms, friendly staff",
		Cons:         "Limited parking",
		Source:       "booking.com",
	}

	hotelID := 123
	expectedID := 456
	createdAt := time.Now()

	mock.ExpectQuery(`INSERT INTO reviews`).
		WithArgs(
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
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "hotel_id", "created_at"}).
			AddRow(expectedID, hotelID, createdAt))

	err = reviewModel.Insert(hotelID, review)

	if err != nil {
		t.Errorf("error was not expected while inserting review: %s", err)
	}

	if review.ID != expectedID {
		t.Errorf("expected ID to be %d, got %d", expectedID, review.ID)
	}

	if review.HotelID != hotelID {
		t.Errorf("expected HotelID to be %d, got %d", hotelID, review.HotelID)
	}

	if review.CreatedAt != createdAt {
		t.Errorf("expected CreatedAt to be %v, got %v", createdAt, review.CreatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_Insert_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	review := &Review{
		AverageScore: 8,
		Name:         "John Doe",
	}

	hotelID := 123

	mock.ExpectQuery(`INSERT INTO reviews`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	err = reviewModel.Insert(hotelID, review)

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

func TestReviewModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	expectedReview := &Review{
		ID:           456,
		HotelID:      123,
		AverageScore: 8,
		Country:      "USA",
		Type:         "Business",
		Name:         "John Doe",
		Date:         "2024-01-15",
		Headline:     "Great stay!",
		Language:     "en",
		Pros:         "Clean rooms, friendly staff",
		Cons:         "Limited parking",
		Source:       "booking.com",
		CreatedAt:    time.Now(),
	}

	mock.ExpectQuery(`SELECT id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE id = \$1 AND hotel_id = \$2`).
		WithArgs(int64(456), int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "hotel_id", "average_score", "country", "type", "name",
			"date", "headline", "language", "pros", "cons", "source", "created_at",
		}).AddRow(
			expectedReview.ID, expectedReview.HotelID, expectedReview.AverageScore,
			expectedReview.Country, expectedReview.Type, expectedReview.Name,
			expectedReview.Date, expectedReview.Headline, expectedReview.Language,
			expectedReview.Pros, expectedReview.Cons, expectedReview.Source,
			expectedReview.CreatedAt,
		))

	review, err := reviewModel.Get(123, 456)

	if err != nil {
		t.Errorf("error was not expected while getting review: %s", err)
	}

	if !reflect.DeepEqual(review, expectedReview) {
		t.Errorf("expected review %+v, got %+v", expectedReview, review)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	mock.ExpectQuery(`SELECT id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE id = \$1 AND hotel_id = \$2`).
		WithArgs(int64(999), int64(123)).
		WillReturnError(sql.ErrNoRows)

	review, err := reviewModel.Get(123, 999)

	if review != nil {
		t.Error("expected review to be nil")
	}

	if err != ErrRecordNotFound {
		t.Errorf("expected error to be %v, got %v", ErrRecordNotFound, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_Get_InvalidID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	testCases := []int64{0, -1, -999}

	for _, id := range testCases {
		t.Run("invalid_id", func(t *testing.T) {
			review, err := reviewModel.Get(123, id)

			if review != nil {
				t.Error("expected review to be nil")
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

func TestReviewModel_Get_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	mock.ExpectQuery(`SELECT id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE id = \$1 AND hotel_id = \$2`).
		WithArgs(int64(456), int64(123)).
		WillReturnError(sql.ErrConnDone)

	review, err := reviewModel.Get(123, 456)

	if review != nil {
		t.Error("expected review to be nil")
	}

	if err != sql.ErrConnDone {
		t.Errorf("expected error to be %v, got %v", sql.ErrConnDone, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafelist: []string{"id", "name"},
	}

	hotelID := int64(123)

	expectedReviews := []*Review{
		{
			ID:           456,
			HotelID:      123,
			AverageScore: 8,
			Country:      "USA",
			Type:         "Business",
			Name:         "John Doe",
			Date:         "2024-01-15",
			Headline:     "Great stay!",
			Language:     "en",
			Pros:         "Clean rooms, friendly staff",
			Cons:         "Limited parking",
			Source:       "booking.com",
			CreatedAt:    time.Now(),
		},
		{
			ID:           457,
			HotelID:      123,
			AverageScore: 9,
			Country:      "Canada",
			Type:         "Leisure",
			Name:         "Jane Smith",
			Date:         "2024-01-16",
			Headline:     "Excellent service!",
			Language:     "en",
			Pros:         "Great location, amazing breakfast",
			Cons:         "WiFi could be better",
			Source:       "expedia.com",
			CreatedAt:    time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"count", "id", "hotel_id", "average_score", "country", "type", "name",
		"date", "headline", "language", "pros", "cons", "source", "created_at",
	})

	for _, review := range expectedReviews {
		rows.AddRow(
			2, review.ID, review.HotelID, review.AverageScore,
			review.Country, review.Type, review.Name, review.Date,
			review.Headline, review.Language, review.Pros, review.Cons,
			review.Source, review.CreatedAt,
		)
	}

	mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE hotel_id = \$1 AND \(fts @@ plainto_tsquery\('simple', \$2\) OR \$2 = ''\) ORDER BY id ASC, id ASC LIMIT \$3 OFFSET \$4`).
		WithArgs(hotelID, "test search", 20, 0).
		WillReturnRows(rows)

	reviews, metadata, err := reviewModel.GetAll(hotelID, "test search", filters)

	if err != nil {
		t.Errorf("error was not expected while getting all reviews: %s", err)
	}

	if len(reviews) != 2 {
		t.Errorf("expected 2 reviews, got %d", len(reviews))
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

	// Verify the content of the first review
	if reviews[0].ID != expectedReviews[0].ID {
		t.Errorf("expected first review ID to be %d, got %d", expectedReviews[0].ID, reviews[0].ID)
	}

	if reviews[0].Name != expectedReviews[0].Name {
		t.Errorf("expected first review Name to be %s, got %s", expectedReviews[0].Name, reviews[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_GetAll_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafelist: []string{"id", "name"},
	}

	hotelID := int64(123)

	rows := sqlmock.NewRows([]string{
		"count", "id", "hotel_id", "average_score", "country", "type", "name",
		"date", "headline", "language", "pros", "cons", "source", "created_at",
	})

	mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE hotel_id = \$1 AND \(fts @@ plainto_tsquery\('simple', \$2\) OR \$2 = ''\) ORDER BY id ASC, id ASC LIMIT \$3 OFFSET \$4`).
		WithArgs(hotelID, "", 20, 0).
		WillReturnRows(rows)

	reviews, metadata, err := reviewModel.GetAll(hotelID, "", filters)

	if err != nil {
		t.Errorf("error was not expected while getting all reviews: %s", err)
	}

	if len(reviews) != 0 {
		t.Errorf("expected 0 reviews, got %d", len(reviews))
	}

	if metadata.TotalRecords != 0 {
		t.Errorf("expected TotalRecords to be 0, got %d", metadata.TotalRecords)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_GetAll_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafelist: []string{"id", "name"},
	}

	hotelID := int64(123)

	mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE hotel_id = \$1 AND \(fts @@ plainto_tsquery\('simple', \$2\) OR \$2 = ''\) ORDER BY id ASC, id ASC LIMIT \$3 OFFSET \$4`).
		WithArgs(hotelID, "", 20, 0).
		WillReturnError(sql.ErrConnDone)

	reviews, metadata, err := reviewModel.GetAll(hotelID, "", filters)

	if err == nil {
		t.Error("expected error, but got none")
	}

	if err != sql.ErrConnDone {
		t.Errorf("expected error to be %v, got %v", sql.ErrConnDone, err)
	}

	if reviews != nil {
		t.Error("expected reviews to be nil")
	}

	if metadata.TotalRecords != 0 {
		t.Errorf("expected empty metadata, got %+v", metadata)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestReviewModel_GetAll_RowScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafelist: []string{"id", "name"},
	}

	hotelID := int64(123)

	// Create a row with invalid data type that will cause scan error
	rows := sqlmock.NewRows([]string{
		"count", "id", "hotel_id", "average_score", "country", "type", "name",
		"date", "headline", "language", "pros", "cons", "source", "created_at",
	}).AddRow(
		"invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid",
		"invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid",
	)

	mock.ExpectQuery(`SELECT count\(\*\) OVER\(\), id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE hotel_id = \$1 AND \(fts @@ plainto_tsquery\('simple', \$2\) OR \$2 = ''\) ORDER BY id ASC, id ASC LIMIT \$3 OFFSET \$4`).
		WithArgs(hotelID, "", 20, 0).
		WillReturnRows(rows)

	reviews, metadata, err := reviewModel.GetAll(hotelID, "", filters)

	if err == nil {
		t.Error("expected error, but got none")
	}

	if reviews != nil {
		t.Error("expected reviews to be nil")
	}

	if metadata.TotalRecords != 0 {
		t.Errorf("expected empty metadata, got %+v", metadata)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewReviewModel(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := NewReviewModel(db)

	if reviewModel == nil {
		t.Error("expected reviewModel to not be nil")
		return
	}

	if reviewModel.DB != db {
		t.Error("expected reviewModel.DB to be the same as the provided db")
	}
}

// Benchmark tests
func BenchmarkReviewModel_Get(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(`SELECT id, hotel_id, average_score, country, type, name, date, headline, language, pros, cons, source, created_at FROM reviews WHERE id = \$1 AND hotel_id = \$2`).
			WithArgs(int64(456), int64(123)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "hotel_id", "average_score", "country", "type", "name",
				"date", "headline", "language", "pros", "cons", "source", "created_at",
			}).AddRow(
				456, 123, 8, "USA", "Business", "John Doe",
				"2024-01-15", "Great stay!", "en", "Clean rooms", "Limited parking",
				"booking.com", time.Now(),
			))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reviewModel.Get(123, 456)
		if err != nil {
			b.Errorf("unexpected error: %s", err)
		}
	}
}

func BenchmarkReviewModel_Insert(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	reviewModel := ReviewModel{DB: db}

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(`INSERT INTO reviews`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "hotel_id", "created_at"}).
				AddRow(456, 123, time.Now()))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		review := &Review{
			AverageScore: 8,
			Country:      "USA",
			Type:         "Business",
			Name:         "John Doe",
			Date:         "2024-01-15",
			Headline:     "Great stay!",
			Language:     "en",
			Pros:         "Clean rooms",
			Cons:         "Limited parking",
			Source:       "booking.com",
		}

		err := reviewModel.Insert(123, review)
		if err != nil {
			b.Errorf("unexpected error: %s", err)
		}
	}
}
