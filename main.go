package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eganbarov/verification_code_service/generator"
	"github.com/eganbarov/verification_code_service/handler"
	"github.com/eganbarov/verification_code_service/lock"
	"github.com/eganbarov/verification_code_service/repository"
	"github.com/eganbarov/verification_code_service/sender"
	"github.com/eganbarov/verification_code_service/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	startServer()
}

func startServer() {
	db := initRedisStorage()
	codeRepository := repository.CodeRepository{Redis: db}
	locker := lock.RedisLocker{Redis: db}
	codeGenerator := generator.CodeGenerator{}
	codeSender := sender.SmsSender{}

	mux := http.NewServeMux()
	mux.Handle(
		"POST /send-code",
		&handler.SendCodeHandler{
			CodeRepository: &codeRepository,
			Locker:         &locker,
			CodeGenerator:  &codeGenerator,
			CodeSender:     &codeSender,
		},
	)
	mux.Handle(
		"POST /validate-code",
		&handler.ValidateCodeHandler{
			CodeRepository: &codeRepository,
			Locker:         &locker,
		},
	)

	listenPort := os.Getenv("LISTEN_PORT")
	fmt.Println("Server starting on :" + listenPort)
	if err := http.ListenAndServe(":"+listenPort, mux); err != nil {
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
