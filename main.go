package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eganbarov/verification_code_service/handler"
	"github.com/eganbarov/verification_code_service/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	startServer()
}

func startServer() {
	db := initRedisStorage()

	mux := http.NewServeMux()
	mux.Handle("POST /send-code", &handler.SendCodeHandler{Redis: db})
	mux.Handle("POST /validate-code", &handler.ValidateCodeHandler{Redis: db})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func initRedisStorage() *redis.Client {
	cnf := storage.Config{
		Addr:        os.Getenv("REDIS_HOST"),
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}
	db, err := storage.NewClient(context.Background(), cnf)
	if err != nil {
		panic(err)
	}

	return db
}
