package countries

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFetchByName_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := &RestCountriesClient{http: &http.Client{Timeout: time.Second}, base: ts.URL}
	_, err := c.FetchByName(context.Background(), "Bharat")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestFetchByName(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Contains(t, r.URL.Path, "name/India")
		resp := []map[string]any{
			{
				"name":    map[string]any{"common": "India"},
				"capital": []string{"Delhi"},
				"currencies": map[string]any{
					"INR": map[string]any{"symbol": "₹"},
				},
				"population": 100,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()
	c := &RestCountriesClient{
		http: &http.Client{Timeout: time.Second},
		base: testServer.URL,
	}
	result, err := c.FetchByName(context.Background(), "India")
	require.NoError(t, err)
	require.Equal(t, "India", result.Name)
	require.Equal(t, "Delhi", result.Capital)
	require.Equal(t, "₹", result.Currency)
	require.EqualValues(t, 100, result.Population)

	restcountry := NewRestCountriesClient(time.Second)
	require.Equal(t, time.Second, restcountry.http.Timeout)
}
func TestFetchByName_BadRequestURL(t *testing.T) {
	c := &RestCountriesClient{
		http: &http.Client{Timeout: time.Second},
		base: "http://wrongurl",
	}
	_, err := c.FetchByName(context.Background(), "India")
	require.Error(t, err)
}
