package data

import (
	"testing"

	"github.com/JLL32/nuitee/internal/validator"
)

func TestFilters_sortColumn(t *testing.T) {
	tests := []struct {
		name        string
		filters     Filters
		expected    string
		shouldPanic bool
	}{
		{
			name: "valid sort column without prefix",
			filters: Filters{
				Sort:         "name",
				SortSafelist: []string{"name", "id", "-name", "-id"},
			},
			expected:    "name",
			shouldPanic: false,
		},
		{
			name: "valid sort column with prefix",
			filters: Filters{
				Sort:         "-name",
				SortSafelist: []string{"name", "id", "-name", "-id"},
			},
			expected:    "name",
			shouldPanic: false,
		},
		{
			name: "invalid sort column",
			filters: Filters{
				Sort:         "invalid",
				SortSafelist: []string{"name", "id"},
			},
			expected:    "",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic but didn't get one")
					}
				}()
			}

			result := tt.filters.sortColumn()

			if !tt.shouldPanic && result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFilters_sortDirection(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		expected string
	}{
		{
			name:     "ascending sort",
			sort:     "name",
			expected: "ASC",
		},
		{
			name:     "descending sort",
			sort:     "-name",
			expected: "DESC",
		},
		{
			name:     "empty sort",
			sort:     "",
			expected: "ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := Filters{Sort: tt.sort}
			result := filters.sortDirection()

			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFilters_limit(t *testing.T) {
	filters := Filters{PageSize: 50}
	expected := 50

	result := filters.limit()

	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestFilters_offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 20,
			expected: 0,
		},
		{
			name:     "second page",
			page:     2,
			pageSize: 20,
			expected: 20,
		},
		{
			name:     "third page with different page size",
			page:     3,
			pageSize: 10,
			expected: 20,
		},
		{
			name:     "large page number",
			page:     100,
			pageSize: 25,
			expected: 2475,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := Filters{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}

			result := filters.offset()

			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestValidateFilters(t *testing.T) {
	tests := []struct {
		name        string
		filters     Filters
		expectValid bool
		errorFields []string
	}{
		{
			name: "valid filters",
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: true,
			errorFields: []string{},
		},
		{
			name: "invalid page - zero",
			filters: Filters{
				Page:         0,
				PageSize:     20,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page"},
		},
		{
			name: "invalid page - negative",
			filters: Filters{
				Page:         -1,
				PageSize:     20,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page"},
		},
		{
			name: "invalid page - too large",
			filters: Filters{
				Page:         10_000_001,
				PageSize:     20,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page"},
		},
		{
			name: "invalid page size - zero",
			filters: Filters{
				Page:         1,
				PageSize:     0,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page_size"},
		},
		{
			name: "invalid page size - negative",
			filters: Filters{
				Page:         1,
				PageSize:     -1,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page_size"},
		},
		{
			name: "invalid page size - too large",
			filters: Filters{
				Page:         1,
				PageSize:     101,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page_size"},
		},
		{
			name: "invalid sort - not in safelist",
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "invalid",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"sort"},
		},
		{
			name: "multiple invalid fields",
			filters: Filters{
				Page:         0,
				PageSize:     101,
				Sort:         "invalid",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: false,
			errorFields: []string{"page", "page_size", "sort"},
		},
		{
			name: "valid sort with descending prefix",
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-name",
				SortSafelist: []string{"name", "id", "-name", "-id"},
			},
			expectValid: true,
			errorFields: []string{},
		},
		{
			name: "boundary values - valid",
			filters: Filters{
				Page:         10_000_000,
				PageSize:     100,
				Sort:         "name",
				SortSafelist: []string{"name", "id"},
			},
			expectValid: true,
			errorFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateFilters(v, tt.filters)

			if tt.expectValid && !v.Valid() {
				t.Errorf("expected valid filters but got errors: %v", v.Errors)
			}

			if !tt.expectValid && v.Valid() {
				t.Error("expected invalid filters but validation passed")
			}

			// Check that specific error fields are present
			for _, field := range tt.errorFields {
				if _, exists := v.Errors[field]; !exists {
					t.Errorf("expected error for field %s but it was not found", field)
				}
			}
		})
	}
}

func TestCalculateMetaData(t *testing.T) {
	tests := []struct {
		name         string
		totalRecords int
		page         int
		pageSize     int
		expected     Metadata
	}{
		{
			name:         "no records",
			totalRecords: 0,
			page:         1,
			pageSize:     20,
			expected:     Metadata{},
		},
		{
			name:         "single page",
			totalRecords: 15,
			page:         1,
			pageSize:     20,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     20,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 15,
			},
		},
		{
			name:         "multiple pages - first page",
			totalRecords: 45,
			page:         1,
			pageSize:     20,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     20,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 45,
			},
		},
		{
			name:         "multiple pages - middle page",
			totalRecords: 100,
			page:         3,
			pageSize:     20,
			expected: Metadata{
				CurrentPage:  3,
				PageSize:     20,
				FirstPage:    1,
				LastPage:     5,
				TotalRecords: 100,
			},
		},
		{
			name:         "multiple pages - last page",
			totalRecords: 100,
			page:         5,
			pageSize:     20,
			expected: Metadata{
				CurrentPage:  5,
				PageSize:     20,
				FirstPage:    1,
				LastPage:     5,
				TotalRecords: 100,
			},
		},
		{
			name:         "exact division",
			totalRecords: 40,
			page:         2,
			pageSize:     20,
			expected: Metadata{
				CurrentPage:  2,
				PageSize:     20,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 40,
			},
		},
		{
			name:         "small page size",
			totalRecords: 23,
			page:         3,
			pageSize:     5,
			expected: Metadata{
				CurrentPage:  3,
				PageSize:     5,
				FirstPage:    1,
				LastPage:     5,
				TotalRecords: 23,
			},
		},
		{
			name:         "large page size",
			totalRecords: 50,
			page:         1,
			pageSize:     100,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMetaData(tt.totalRecords, tt.page, tt.pageSize)

			if result.CurrentPage != tt.expected.CurrentPage {
				t.Errorf("CurrentPage: expected %d, got %d", tt.expected.CurrentPage, result.CurrentPage)
			}

			if result.PageSize != tt.expected.PageSize {
				t.Errorf("PageSize: expected %d, got %d", tt.expected.PageSize, result.PageSize)
			}

			if result.FirstPage != tt.expected.FirstPage {
				t.Errorf("FirstPage: expected %d, got %d", tt.expected.FirstPage, result.FirstPage)
			}

			if result.LastPage != tt.expected.LastPage {
				t.Errorf("LastPage: expected %d, got %d", tt.expected.LastPage, result.LastPage)
			}

			if result.TotalRecords != tt.expected.TotalRecords {
				t.Errorf("TotalRecords: expected %d, got %d", tt.expected.TotalRecords, result.TotalRecords)
			}
		})
	}
}

// Benchmark tests
func BenchmarkFilters_sortColumn(b *testing.B) {
	filters := Filters{
		Sort:         "name",
		SortSafelist: []string{"name", "id", "-name", "-id"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = filters.sortColumn()
	}
}

func BenchmarkFilters_sortDirection(b *testing.B) {
	filters := Filters{Sort: "-name"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = filters.sortDirection()
	}
}

func BenchmarkFilters_offset(b *testing.B) {
	filters := Filters{
		Page:     10,
		PageSize: 20,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = filters.offset()
	}
}

func BenchmarkCalculateMetaData(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = calculateMetaData(1000, 5, 20)
	}
}

func BenchmarkValidateFilters(b *testing.B) {
	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "name",
		SortSafelist: []string{"name", "id", "-name", "-id"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := validator.New()
		ValidateFilters(v, filters)
	}
}
