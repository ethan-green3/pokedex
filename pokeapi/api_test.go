package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocationAreas_Success(t *testing.T) {
	mockResponse := `{
		"results": [
			{
				"name": "canalave-city",
				"url": "https://pokeapi.co/api/v2/location/1/"
			},
			{
				"name": "eterna-city",
				"url": "https://pokeapi.co/api/v2/location/2/"
			}
		],
		"next": "https://pokeapi.co/api/v2/location/?offset=20",
		"previous": null
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	resp, err := GetLocationAreas(server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(resp.Results))
	}

	if resp.Results[0].Name != "canalave-city" {
		t.Errorf("expected canalave-city, got %s", resp.Results[0].Name)
	}

	if resp.Next == nil {
		t.Errorf("expected next page URL but got nil")
	}
}

func TestGetLocationAreas_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	_, err := GetLocationAreas(server.URL)

	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}

func TestGetLocationAreas_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := GetLocationAreas(server.URL)
	if err == nil {
		t.Fatalf("expected error from HTTP 500 response, got nil")
	}
}
