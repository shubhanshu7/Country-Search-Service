package service

import (
	"context"
	"countrySearchService/internal/cache"
	"countrySearchService/internal/countries"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeClient struct {
	delay   time.Duration
	resp    countries.Country
	err     error
	callsMu sync.Mutex
	calls   int
}

func (f *fakeClient) FetchByName(ctx context.Context, name string) (countries.Country, error) {
	f.callsMu.Lock()
	f.calls++
	f.callsMu.Unlock()
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return countries.Country{}, ctx.Err()
		}
	}
	if f.err != nil {
		return countries.Country{}, f.err
	}
	return f.resp, nil
}

func TestService_CacheHitMiss(t *testing.T) {
	mem := cache.NewMemoryCache()
	fc := &fakeClient{resp: countries.Country{Name: "India"}}
	s := NewCountryService(mem, fc)

	// set
	got1, err := s.Search(context.Background(), "India")
	require.NoError(t, err)
	require.Equal(t, "India", got1.Name)

	//  no new fetch
	got2, err := s.Search(context.Background(), "India")
	require.NoError(t, err)
	require.Equal(t, "India", got2.Name)

	fc.callsMu.Lock()
	defer fc.callsMu.Unlock()
	require.Equal(t, 1, fc.calls)
}

func TestService_ValidateName(t *testing.T) {
	mem := cache.NewMemoryCache()
	s := NewCountryService(mem, &fakeClient{})
	_, err := s.Search(context.Background(), "")
	require.Error(t, err)
}
