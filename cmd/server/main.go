package main

import (
	"context"
	"countrySearchService/internal/httpapi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//cache
	//client timeout
	//service layer
	//api

	mux := http.NewServeMux()
	httpapi.Routes(mux)

	srv := &http.Server{
		Addr:        ":8000",
		Handler:     mux,
		IdleTimeout: 30 * time.Second,
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
