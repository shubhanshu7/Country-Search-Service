package service

import (
	"context"
	"countrySearchService/internal/cache"
	"countrySearchService/internal/countries"
	"errors"
	"fmt"
)

type CountryService struct {
	cache  cache.Cache
	client countries.Client
}

func NewCountryService(c cache.Cache, client countries.Client) *CountryService {
	return &CountryService{cache: c, client: client}
}
func (s *CountryService) Search(ctx context.Context, name string) (countries.Country, error) {
	if name == "" {
		return countries.Country{}, errors.New("name is required")
	}
	if value, ok := s.cache.Get(name); ok {
		fmt.Println("=======getting value from cache================")
		return value.(countries.Country), nil
	}

	country, err := s.client.FetchByName(ctx, name)
	if err != nil {
		return countries.Country{}, err
	}
	s.cache.Set(name, country)
	return country, nil
}
