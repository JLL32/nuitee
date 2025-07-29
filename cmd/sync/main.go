package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
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
	HotelID      int     `json:"hotel_id"`
	MainImageTh  string  `json:"main_image_th"`
	HotelName    string  `json:"hotel_name"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Address      Address `json:"address"`
	Stars        int     `json:"stars"`
	Rating       float64 `json:"rating"`
	ReviewCount  int     `json:"review_count"`
	ChildAllowed bool    `json:"child_allowed"`
	PetsAllowed  bool    `json:"pets_allowed"`
	Description  string  `json:"description"`
}

type Review struct {
	AverageScore int    `json:"average_score"`
	Country      string `json:"country"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Date         string `json:"date"`
	Headline     string `json:"headline"`
	Language     string `json:"language"`
	Pros         string `json:"pros"`
	Cons         string `json:"cons"`
}

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

	if inputFile == "" || dsn == "" {
		flag.Usage()
		return
	}

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	ids := strings.Split(string(inputData), ",")
	for i := range ids {
		ids[i] = strings.TrimSpace(ids[i])
	}

	for _, id := range ids {
		var hotel Hotel
		err := fetchJson(fmt.Sprintf("%s/v3.0/property/%s", apiUrl, id), map[string]string{"x-api-key": apiKey}, &hotel)
		if err != nil {
			fmt.Printf("Error fetching hotel data for ID %s: %v\n", id, err)
			continue
		}

		var reviews []Review
		err = fetchJson(fmt.Sprintf("%s/v3.0/property/reviews/%s/1000000", apiUrl, id), map[string]string{"x-api-key": apiKey}, &reviews)
		if err != nil {
			fmt.Printf("Error fetching review data for ID %s: %v\n", id, err)
			continue
		}

		fmt.Printf("%+v\n", hotel)
		fmt.Printf("%+v\n", reviews)
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
