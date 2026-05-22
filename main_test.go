package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

// TestCafeCount() and TestCafeSearch() have its own values for testing ecah city

func TestCafeCount(t *testing.T) {
	// defining struct for using in map
	type results struct {
		count int
		want  int
	}
	// defining map for using in testing get-requests for both cities
	requests := map[string][]results{
		"moscow": {
			{count: 0, want: 0},
			{count: 1, want: 1},
			{count: 2, want: 2},
			{count: 100, want: 5}, // must show all restaurants, 5 for Moscow
		},
		"tula": {
			{count: 0, want: 0},
			{count: 1, want: 1},
			{count: 2, want: 2},
			{count: 100, want: 3}, // must show all restaurants, 3 for Tula
		},
	}
	// testing mainHandle in main.go
	handler := http.HandlerFunc(mainHandle)

	// comparing responses for each city with map values
	// all get-requests are without "search" parameter
	for city, r := range requests {
		for _, request := range r {
			// rendering URL with unique parameters
			url := fmt.Sprintf("/cafe?city=%s&count=%d", city, request.count)
			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)

			handler.ServeHTTP(response, req)

			result := response.Body.String()
			slice := strings.Split(result, `,`)
			// exception for the "zero" result, cause length of empty slice won't be 0
			if result == "" {
				assert.Equal(t, request.want, 0)
			} else {
				assert.Equal(t, request.want, len(slice))
			}
		}
	}
}

func TestCafeSearch(t *testing.T) {
	// defining struct for using in map
	type results struct {
		search    string
		wantCount int
	}
	// defining map for using in testing get-requests for both cities
	requests := map[string][]results{
		"moscow": {
			{search: "", wantCount: 5},
			{search: "фасоль", wantCount: 0},
			{search: "кофе", wantCount: 2},
			{search: "вилка", wantCount: 1},
		},
		"tula": {
			{search: "", wantCount: 3},
			{search: "мир", wantCount: 1},
			{search: "завтрак", wantCount: 1},
			{search: "за", wantCount: 2}, // test for the part of word
		},
	}
	// testing mainhandle in main.go
	handler := http.HandlerFunc(mainHandle)

	// comparing responses for each city with map values
	// all requests are without "count" parameter
	for city, r := range requests {
		for _, request := range r {
			// rendering URL with unique parameters
			url := fmt.Sprintf("/cafe?city=%s&search=%s", city, request.search)
			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)

			handler.ServeHTTP(response, req)

			result := response.Body.String()
			slice := strings.Split(result, `,`)
			// exception for the "zero" result, cause length of empty slice won't be 0
			if result == "" {
				assert.Equal(t, request.wantCount, 0)
			} else {
				assert.Equal(t, request.wantCount, len(slice))
			}
		}
	}
}
