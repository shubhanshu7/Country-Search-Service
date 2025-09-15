package main

import (
	"context"
	"countrySearchService/internal/cache"
	"countrySearchService/internal/countries"
	"countrySearchService/internal/httpapi"
	"countrySearchService/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//cache
	mem := cache.NewMemoryCache()

	//client timeout
	client := countries.NewRestCountriesClient(5 * time.Second)
	//service layer
	svc := service.NewCountryService(mem, client)
	//api
	api := httpapi.NewServer(svc)

	mux := http.NewServeMux()
	api.Routes(mux)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
	log.Println("server stopped")
}
