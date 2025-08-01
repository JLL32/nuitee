package main

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JLL32/nuitee/internal/data"
	"github.com/julienschmidt/httprouter"
)

func newTestApplication(t *testing.T) (*application, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	cfg := config{
		port: 4000,
		env:  "test",
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}{
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  15 * time.Minute,
		},
		limiter: struct {
			rps     float64
			burst   int
			enabled bool
		}{
			rps:     2,
			burst:   4,
			enabled: false, // Disable rate limiting for tests
		},
		openAIkey: "test-key",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	return app, mock, func() {
		db.Close()
	}
}

func newTestApplicationForBenchmark(b *testing.B) (*application, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	cfg := config{
		port: 4000,
		env:  "test",
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}{
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  15 * time.Minute,
		},
		limiter: struct {
			rps     float64
			burst   int
			enabled bool
		}{
			rps:     2,
			burst:   4,
			enabled: false, // Disable rate limiting for tests
		},
		openAIkey: "test-key",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	return app, mock, func() {
		db.Close()
	}
}

// testRoutes returns routes without metrics middleware to avoid expvar conflicts
func (app *application) testRoutes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/hotels", app.listHotelsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID", app.getHotelHandler)

	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews", app.listReviewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews/:reviewID", app.getReviewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews/:reviewID/summary", app.getReviewSummaryHandler)

	return app.recoverPanic(app.rateLimit(router))
}