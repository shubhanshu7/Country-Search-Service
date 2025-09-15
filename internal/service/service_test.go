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

	got1, err := s.Search(context.Background(), "India")
	require.NoError(t, err)
	require.Equal(t, "India", got1.Name)

	got2, err := s.Search(context.Background(), "India")
	require.NoError(t, err)
	require.Equal(t, "India", got2.Name)

	fc.callsMu.Lock()
	defer fc.callsMu.Unlock()
	require.Equal(t, 1, fc.calls)
}

func TestService_ConcurrentRequests(t *testing.T) {
	mem := cache.NewMemoryCache()
	fc := &fakeClient{
		delay: 20 * time.Millisecond,
		resp:  countries.Country{Name: "India"},
	}
	s := NewCountryService(mem, fc)

	var wg sync.WaitGroup
	results := make(chan error, 50)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.Search(context.Background(), "India")
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}

	fc.callsMu.Lock()
	calls := fc.calls
	fc.callsMu.Unlock()
	require.Less(t, calls, 51)
}

func TestService_ValidateName(t *testing.T) {
	mem := cache.NewMemoryCache()
	s := NewCountryService(mem, &fakeClient{})
	_, err := s.Search(context.Background(), "")
	require.Error(t, err)
}
