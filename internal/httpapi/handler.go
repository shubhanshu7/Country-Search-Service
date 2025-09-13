package httpapi

import (
	"context"
	"countrySearchService/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	svc *service.CountryService
}

func NewServer(svc *service.CountryService) *Server {
	return &Server{svc: svc}
}
func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/countries/search", s.handleSearch)
}
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=============handler start=============")
	name := r.URL.Query().Get("name")
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	country, err := s.svc.Search(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("=====country====", country)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
