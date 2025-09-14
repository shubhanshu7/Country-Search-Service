package httpapi

import (
	"context"
	"countrySearchService/internal/cache"
	"countrySearchService/internal/countries"
	"countrySearchService/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type okClient struct{ c countries.Country }

func (o okClient) FetchByName(ctx context.Context, name string) (countries.Country, error) {
	return o.c, nil
}

func TestSearchHandler_OK(t *testing.T) {
	mem := cache.NewMemoryCache()
	svc := service.NewCountryService(mem, okClient{c: countries.Country{
		Name: "India", Capital: "New Delhi", Currency: "₹", Population: 1380004385,
	}})
	s := NewServer(svc)
	mux := http.NewServeMux()
	s.Routes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=India", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got countries.Country
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, "India", got.Name)
}

func TestSearchHandler_BadRequest(t *testing.T) {
	mem := cache.NewMemoryCache()
	svc := service.NewCountryService(mem, okClient{})
	s := NewServer(svc)
	mux := http.NewServeMux()
	s.Routes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
