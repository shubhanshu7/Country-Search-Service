package countries

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Country struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int64  `json:"population"`
}

type RestCountriesClient struct {
	http *http.Client
	base string
}

func NewRestCountriesClient(timeout time.Duration) *RestCountriesClient {
	return &RestCountriesClient{
		http: &http.Client{Timeout: timeout},
		base: "https://restcountries.com/v3.1",
	}
}

type Client interface {
	FetchByName(ctx context.Context, name string) (Country, error)
}

var ErrNotFound = errors.New("country not found")

type rcName struct {
	Common string `json:"common"`
}
type rcCurrencies map[string]struct {
	Symbol string `json:"symbol"`
}
type rcCountry struct {
	Name       rcName       `json:"name"`
	Capital    []string     `json:"capital"`
	Currencies rcCurrencies `json:"currencies"`
	Population int64        `json:"population"`
}

func (c *RestCountriesClient) FetchByName(ctx context.Context, name string) (Country, error) {
	url := fmt.Sprintf("%s/name/%s?fields=name,capital,currencies,population", c.base, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Country{}, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return Country{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Country{}, fmt.Errorf("error fetching %s named country", name)
	}
	var resultCountry []rcCountry
	if err := json.NewDecoder(resp.Body).Decode(&resultCountry); err != nil {
		return Country{}, err
	}
	if len(resultCountry) == 0 {
		return Country{}, ErrNotFound
	}
	id := 0
	for i := range resultCountry {
		if strings.EqualFold(resultCountry[i].Name.Common, name) {
			id = i
			break
		}
	}
	finalCountry := resultCountry[id]
	fmt.Println("=============final output=====", finalCountry)
	currency := ""
	for _, curr := range finalCountry.Currencies {
		currency = curr.Symbol
		break
	}
	capital := ""
	if len(finalCountry.Capital) > 0 {
		capital = finalCountry.Capital[0]
	}
	return Country{
		Name:       finalCountry.Name.Common,
		Capital:    capital,
		Currency:   currency,
		Population: finalCountry.Population,
	}, nil
}
