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
	_ "github.com/lib/pq"
)

/*
hotel:
 * {
   "hotel_id": 1022380,
   "main_image_th": "https://static.cupid.travel/hotels/thumbnail/480020386.jpg",
   "hotel_name": "London Marriott Hotel Marble Arch",
   "phone": "",
   "email": "",
   "address": {
     "address": "134 George Street",
     "city": "London",
     "state": "Greater London",
     "country": "gb",
     "postal_code": "W1H 5DN"
   },
   "stars": 4,
   "rating": 7.5,
   "review_count": 1301,
   "child_allowed": true,
   "pets_allowed": false,
   "description": "<p><strong>Modern Fitness Centre and Stylish Rooms</strong><br>The Marriott Hotel Marble Arch offers a modern fitness centre for guests to stay active during their stay. The stylish rooms feature luxurious cotton bedding, Apple TV technology, and a spacious work desk for ultimate comfort and convenience.</p><p><strong>Contemporary Restaurant - The Pickled Hen</strong><br>Indulge in traditional British food and drinks at The Pickled Hen, a modern gastropub located within the hotel. Enjoy a culinary experience without having to leave the premises.</p><p><strong>Prime Location near Oxford Street and Hyde Park</strong><br>Situated just a 5-minute walk from Oxford Street and 500 metres from Hyde Park, Marriott Hotel Marble Arch provides easy access to London's top attractions. Explore the vibrant shopping of Knightsbridge and the lively atmosphere of Soho within walking distance.</p><p>Experience luxury and convenience at Marriott Hotel Marble Arch. Book your stay now for an unforgettable London getaway.</p>",
 }

 review:
 {
   "average_score": 3,
   "country": "au",
   "type": "young couple",
   "name": "Susan",
   "date": "2024-11-10 20:32:40",
   "headline": "Poor",
   "language": "en",
   "pros": "",
   "cons": "",
   "source": "Nuitee"
 },
*/

func main() {
	var (
		inputFile string
		dsn       string
		apiKey    string
		apiUrl    string
	)

	flag.StringVar(&inputFile, "input", "", "Input file path")
	flag.StringVar(&dsn, "db-dsn", "", "Database connection string")
	flag.StringVar(&apiKey, "api-key", "", "API key for authentication")
	flag.StringVar(&apiUrl, "api-url", "", "API URL for fetching data")
	flag.Parse()

	if inputFile == "" || dsn == "" || apiKey == "" || apiUrl == "" {
		flag.Usage()
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	ids := strings.Split(string(inputData), ",")
	for i := range ids {
		ids[i] = strings.TrimSpace(ids[i])
	}

	// TODO: add concurrency
	for _, id := range ids {
		var hotel data.Hotel
		for range 3 {
			err = fetchJson(fmt.Sprintf("%s/v3.0/property/%s", apiUrl, id), map[string]string{"x-api-key": apiKey}, &hotel)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Printf("Error fetching hotel data for ID %s: %v\n", id, err)
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
			fmt.Printf("Error fetching review data for ID %s: %v\n", id, err)
			continue
		}

		db, err := openDB(dsn)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer db.Close()

		// TODO: if we're syncing stuff we should turn this into an upsert function
		models := data.NewModels(db)
		err = models.Hotels.Insert(&hotel)
		if err != nil {
			fmt.Printf("Error inserting hotel data for ID %s: %v\n", id, err)
			continue
		}

		for _, review := range reviews {
			err = models.Reviews.Insert(hotel.HotelID, &review)
			if err != nil {
				fmt.Printf("Error inserting review data for ID %s: %v\n", id, err)
				continue
			}
		}
	}

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
