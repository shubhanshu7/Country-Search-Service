package httpapi

import (
	"context"
	"countrySearchService/internal/cache"
	"countrySearchService/internal/countries"
	"countrySearchService/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type okClient struct{ c countries.Country }

func (o okClient) FetchByName(ctx context.Context, name string) (countries.Country, error) {
	return o.c, nil
}

type boomWriter struct {
	hdr  http.Header
	code int
}

func (b *boomWriter) Header() http.Header {
	if b.hdr == nil {
		b.hdr = http.Header{}
	}
	return b.hdr
}

func (b *boomWriter) WriteHeader(statusCode int) {
	b.code = statusCode
}

func (b *boomWriter) Write(p []byte) (int, error) {
	return 0, errors.New("boom write")
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
func TestSearchHandler_JSONEncodeError(t *testing.T) {
	mem := cache.NewMemoryCache()
	ok := okClient{c: countries.Country{
		Name: "India", Capital: "Delhi", Currency: "₹", Population: 10,
	}}
	svc := service.NewCountryService(mem, ok)
	s := NewServer(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=India", nil)
	bw := &boomWriter{}

	s.handleSearch(bw, req)

	require.Equal(t, http.StatusInternalServerError, bw.code)
}
