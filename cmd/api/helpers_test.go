package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/JLL32/nuitee/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func TestReadIDParam(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	tests := []struct {
		name        string
		paramName   string
		paramValue  string
		expectedID  int64
		expectError bool
	}{
		{
			name:        "valid positive ID",
			paramName:   "id",
			paramValue:  "123",
			expectedID:  123,
			expectError: false,
		},
		{
			name:        "valid large ID",
			paramName:   "hotelID",
			paramValue:  "999999",
			expectedID:  999999,
			expectError: false,
		},
		{
			name:        "invalid zero ID",
			paramName:   "id",
			paramValue:  "0",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "invalid negative ID",
			paramName:   "id",
			paramValue:  "-1",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "invalid non-numeric ID",
			paramName:   "id",
			paramValue:  "abc",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "invalid empty ID",
			paramName:   "id",
			paramValue:  "",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "invalid float ID",
			paramName:   "id",
			paramValue:  "123.45",
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the parameter
			req, err := http.NewRequest("GET", "/test/"+tt.paramValue, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create router context with parameters
			router := httprouter.New()
			router.HandlerFunc(http.MethodGet, "/test/:"+tt.paramName, func(w http.ResponseWriter, r *http.Request) {
				id, err := app.readIDParam(r, tt.paramName)

				if tt.expectError {
					if err == nil {
						t.Errorf("expected error but got none")
					}
				} else {
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
					if id != tt.expectedID {
						t.Errorf("expected ID %d, got %d", tt.expectedID, id)
					}
				}
			})

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
		})
	}
}

func TestWriteJSON(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	tests := []struct {
		name           string
		status         int
		data           envelope
		headers        http.Header
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "simple JSON response",
			status: http.StatusOK,
			data: envelope{
				"message": "hello world",
				"status":  "success",
			},
			headers:        nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if response["message"] != "hello world" {
					t.Errorf("expected message 'hello world', got %v", response["message"])
				}

				if response["status"] != "success" {
					t.Errorf("expected status 'success', got %v", response["status"])
				}

				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected Content-Type 'application/json', got %s", contentType)
				}
			},
		},
		{
			name:   "JSON response with custom headers",
			status: http.StatusCreated,
			data: envelope{
				"id": 123,
			},
			headers: http.Header{
				"X-Custom-Header": []string{"custom-value"},
				"Location":        []string{"/v1/resource/123"},
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				if rr.Header().Get("X-Custom-Header") != "custom-value" {
					t.Errorf("expected X-Custom-Header 'custom-value', got %s", rr.Header().Get("X-Custom-Header"))
				}

				if rr.Header().Get("Location") != "/v1/resource/123" {
					t.Errorf("expected Location '/v1/resource/123', got %s", rr.Header().Get("Location"))
				}
			},
		},
		{
			name:   "empty data",
			status: http.StatusNoContent,
			data:   envelope{},
			headers: nil,
			expectedStatus: http.StatusNoContent,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if len(response) != 0 {
					t.Errorf("expected empty response, got %v", response)
				}
			},
		},
		{
			name:   "complex nested data",
			status: http.StatusOK,
			data: envelope{
				"user": map[string]interface{}{
					"id":   123,
					"name": "John Doe",
					"tags": []string{"admin", "user"},
				},
				"metadata": map[string]interface{}{
					"total": 1,
					"page":  1,
				},
			},
			headers:        nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				user, ok := response["user"].(map[string]interface{})
				if !ok {
					t.Error("expected user to be an object")
				}

				if user["name"] != "John Doe" {
					t.Errorf("expected user name 'John Doe', got %v", user["name"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			err := app.writeJSON(rr, tt.status, tt.data, tt.headers)
			if err != nil {
				t.Fatalf("writeJSON returned error: %v", err)
			}

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}

			tt.checkResponse(t, rr)
		})
	}
}

func TestReadJSON(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	tests := []struct {
		name        string
		body        string
		target      interface{}
		expectError bool
		errorCheck  func(error) bool
	}{
		{
			name: "valid JSON",
			body: `{"name": "John Doe", "age": 30}`,
			target: &struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			expectError: false,
		},
		{
			name:        "invalid JSON syntax",
			body:        `{"name": "John Doe", "age":}`,
			target:      &map[string]interface{}{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "badly-formed JSON")
			},
		},
		{
			name:        "empty body",
			body:        "",
			target:      &map[string]interface{}{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "body must not be empty")
			},
		},
		{
			name:        "unknown field",
			body:        `{"name": "John", "unknown_field": "value"}`,
			target:      &struct {
				Name string `json:"name"`
			}{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "unknown key")
			},
		},
		{
			name: "incorrect type",
			body: `{"name": 123}`,
			target: &struct {
				Name string `json:"name"`
			}{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "incorrect JSON type")
			},
		},
		{
			name:        "multiple JSON values",
			body:        `{"name": "John"}{"age": 30}`,
			target:      &map[string]interface{}{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "single JSON value")
			},
		},
		{
			name:        "unexpected EOF",
			body:        `{"name": "John"`,
			target:      &map[string]interface{}{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "badly-formed JSON")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/test", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			err = app.readJSON(rr, req, tt.target)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorCheck != nil && !tt.errorCheck(err) {
					t.Errorf("error check failed for error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestReadJSON_MaxBytes(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	// Create a body that's larger than 1MB
	largeBody := strings.Repeat("a", 1_048_577) // 1MB + 1 byte
	jsonBody := `{"data": "` + largeBody + `"}`

	req, err := http.NewRequest("POST", "/test", strings.NewReader(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	target := &map[string]interface{}{}

	err = app.readJSON(rr, req, target)

	if err == nil {
		t.Error("expected error for oversized body")
	}

	if !strings.Contains(err.Error(), "must not be larger than") {
		t.Errorf("expected max bytes error, got: %v", err)
	}
}

func TestReadString(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	tests := []struct {
		name         string
		queryValues  url.Values
		key          string
		defaultValue string
		expected     string
	}{
		{
			name: "existing key",
			queryValues: url.Values{
				"search": []string{"hotels"},
			},
			key:          "search",
			defaultValue: "default",
			expected:     "hotels",
		},
		{
			name: "missing key returns default",
			queryValues: url.Values{
				"other": []string{"value"},
			},
			key:          "search",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "empty query values",
			queryValues:  url.Values{},
			key:          "search",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name: "empty string value",
			queryValues: url.Values{
				"search": []string{""},
			},
			key:          "search",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name: "multiple values returns first",
			queryValues: url.Values{
				"search": []string{"first", "second"},
			},
			key:          "search",
			defaultValue: "default",
			expected:     "first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.readString(tt.queryValues, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	tests := []struct {
		name         string
		queryValues  url.Values
		key          string
		defaultValue []string
		expected     []string
	}{
		{
			name: "single value",
			queryValues: url.Values{
				"tags": []string{"hotel"},
			},
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"hotel"},
		},
		{
			name: "multiple values",
			queryValues: url.Values{
				"tags": []string{"hotel,resort,spa"},
			},
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"hotel", "resort", "spa"},
		},
		{
			name: "missing key returns default",
			queryValues: url.Values{
				"other": []string{"value"},
			},
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
		{
			name:         "empty query values",
			queryValues:  url.Values{},
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
		{
			name: "empty string value",
			queryValues: url.Values{
				"tags": []string{""},
			},
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
		{
			name: "values with spaces",
			queryValues: url.Values{
				"tags": []string{"hotel, resort , spa"},
			},
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"hotel", " resort ", " spa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.readCSV(tt.queryValues, tt.key, tt.defaultValue)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %s at index %d, got %s", tt.expected[i], i, v)
				}
			}
		})
	}
}

func TestReadInt(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	tests := []struct {
		name         string
		queryValues  url.Values
		key          string
		defaultValue int
		expected     int
		expectError  bool
	}{
		{
			name: "valid integer",
			queryValues: url.Values{
				"page": []string{"5"},
			},
			key:          "page",
			defaultValue: 1,
			expected:     5,
			expectError:  false,
		},
		{
			name: "missing key returns default",
			queryValues: url.Values{
				"other": []string{"value"},
			},
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  false,
		},
		{
			name:         "empty query values",
			queryValues:  url.Values{},
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  false,
		},
		{
			name: "empty string value",
			queryValues: url.Values{
				"page": []string{""},
			},
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  false,
		},
		{
			name: "invalid integer",
			queryValues: url.Values{
				"page": []string{"abc"},
			},
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  true,
		},
		{
			name: "float value",
			queryValues: url.Values{
				"page": []string{"5.5"},
			},
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  true,
		},
		{
			name: "zero value",
			queryValues: url.Values{
				"page": []string{"0"},
			},
			key:          "page",
			defaultValue: 1,
			expected:     0,
			expectError:  false,
		},
		{
			name: "negative value",
			queryValues: url.Values{
				"page": []string{"-5"},
			},
			key:          "page",
			defaultValue: 1,
			expected:     -5,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			result := app.readInt(tt.queryValues, tt.key, tt.defaultValue, v)

			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}

			if tt.expectError && v.Valid() {
				t.Error("expected validation error but validator is valid")
			}

			if !tt.expectError && !v.Valid() {
				t.Errorf("unexpected validation error: %v", v.Errors)
			}
		})
	}
}

func TestBackground(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	// Test that background function executes
	executed := make(chan bool, 1)

	app.background(func() {
		executed <- true
	})

	// Wait for the background function to execute
	select {
	case <-executed:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("background function did not execute within timeout")
	}

	// Wait for background goroutines to finish
	app.wg.Wait()
}

func TestBackground_Panic(t *testing.T) {
	app, _, cleanup := newTestApplication(t)
	defer cleanup()

	// Test that background function handles panics gracefully
	executed := make(chan bool, 1)

	app.background(func() {
		defer func() {
			executed <- true
		}()
		panic("test panic")
	})

	// Wait for the background function to execute and handle panic
	select {
	case <-executed:
		// Success - panic was handled gracefully
	case <-time.After(100 * time.Millisecond):
		t.Error("background function did not execute within timeout")
	}

	// Wait for background goroutines to finish
	app.wg.Wait()
}



// Benchmark tests
func BenchmarkWriteJSON(b *testing.B) {
	app, _, cleanup := newTestApplicationForBenchmark(b)
	defer cleanup()

	data := envelope{
		"message": "hello world",
		"status":  "success",
		"data": map[string]interface{}{
			"id":   123,
			"name": "test",
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		err := app.writeJSON(rr, http.StatusOK, data, nil)
		if err != nil {
			b.Fatalf("writeJSON error: %v", err)
		}
	}
}

func BenchmarkReadJSON(b *testing.B) {
	app, _, cleanup := newTestApplicationForBenchmark(b)
	defer cleanup()

	jsonBody := `{"name": "John Doe", "age": 30, "email": "john@example.com"}`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		target := &struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Email string `json:"email"`
		}{}

		err := app.readJSON(rr, req, target)
		if err != nil {
			b.Fatalf("readJSON error: %v", err)
		}
	}
}

func BenchmarkReadInt(b *testing.B) {
	app, _, cleanup := newTestApplicationForBenchmark(b)
	defer cleanup()

	qs := url.Values{
		"page": []string{"5"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := validator.New()
		_ = app.readInt(qs, "page", 1, v)
	}
}