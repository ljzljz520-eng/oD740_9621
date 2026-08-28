package main

import (
	"example.com/noise-registry/internal/api"
	"example.com/noise-registry/internal/service"
	"example.com/noise-registry/internal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	path := os.Getenv("NOISE_DB")
	if path == "" {
		path = "noise.db"
	}
	db, e := store.Open(path)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	svc := service.New(db)
	log.Println("noise registry listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.New(svc).Handler()))
}
