package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JLL32/nuitee/internal/data"
	"github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
)

func main() {
	var (
		inputFile string
		dsn       string
		apiKey    string
		apiUrl    string
		interval  int
	)

	flag.StringVar(&inputFile, "input", "", "Input file path")
	flag.StringVar(&dsn, "db-dsn", "", "Database connection string")
	flag.StringVar(&apiKey, "api-key", "", "API key for authentication")
	flag.StringVar(&apiUrl, "api-url", "", "API URL for fetching data")
	flag.IntVar(&interval, "interval", 3, "Interval in minutes")
	flag.Parse()

	if inputFile == "" || dsn == "" || apiKey == "" || apiUrl == "" {
		flag.Usage()
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	ids := strings.Split(string(inputData), ",")
	for i := range ids {
		ids[i] = strings.TrimSpace(ids[i])
	}

	db, err := openDB(dsn)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	defer db.Close()
	slog.Info("Database connection established")

	s := gocron.NewScheduler(time.UTC)
	s.Every(interval).Minutes().Do(func() {
		slog.Info("Starting sync")
		for _, id := range ids {
			var hotel data.Hotel
			for range 3 {
				err = fetchJson(fmt.Sprintf("%s/v3.0/property/%s", apiUrl, id), map[string]string{"x-api-key": apiKey}, &hotel)
				if err == nil {
					break
				}
			}
			if err != nil {
				slog.Error(fmt.Sprintf("Error fetching hotel data for ID %s: %v", id, err))
				continue
			}

			var reviews []data.Review
			for range 3 {
				err = fetchJson(fmt.Sprintf("%s/v3.0/property/reviews/%s/1000000", apiUrl, id), map[string]string{"x-api-key": apiKey}, &reviews)
				if err == nil {
					break
				}
			}
			if err != nil {
				slog.Error(fmt.Sprintf("Error fetching review data for ID %s: %v", id, err))
				continue
			}

			models := data.NewModels(db)
			err = models.Hotels.Upsert(&hotel)
			if err != nil {
				slog.Error(fmt.Sprintf("Error inserting hotel data for ID %s: %v", id, err))
				continue
			}

			for _, review := range reviews {
				err = models.Reviews.Upsert(hotel.HotelID, &review)
				if err != nil {
					slog.Error(fmt.Sprintf("Error inserting review data for ID %s: %v", id, err))
					continue
				}
			}
		}
	})
	s.StartBlocking()

}

func fetchJson(url string, headers map[string]string, output any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(output)
	if err != nil {
		return err
	}

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
