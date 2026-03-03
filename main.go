package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eganbarov/verification_code_service/config"
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
	appConfig := initAppConfig()
	codeRepository := repository.CodeRepository{Redis: db, AppConfig: appConfig}
	locker := lock.RedisLocker{Redis: db, AppConfig: appConfig}
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
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbNumber = storage.DB
	}
	maxRetries, err := strconv.Atoi(os.Getenv("REDIS_MAX_RETRIES"))
	if err != nil {
		maxRetries = storage.MAX_RETRIES
	}
	dialTimeout, err := strconv.Atoi(os.Getenv("REDIS_DIAL_TIMEOUT"))
	if err != nil {
		dialTimeout = storage.DIAL_TIMEOUT
	}
	timeout, err := strconv.Atoi(os.Getenv("REDIS_TIMEOUT"))
	if err != nil {
		timeout = storage.TIMEOUT
	}
	cnf := storage.Config{
		Addr:        os.Getenv("REDIS_HOST"),
		DB:          dbNumber,
		MaxRetries:  maxRetries,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
		Timeout:     time.Duration(timeout) * time.Second,
	}

	db, err := storage.NewClient(context.Background(), cnf)
	if err != nil {
		panic(err)
	}

	return db
}

func initAppConfig() *config.AppConfig {
	codeTtl, err := strconv.Atoi(os.Getenv("CODE_TTL"))
	if err != nil {
		codeTtl = config.CODE_TTL
	}
	repeatSentCodeTtl, err := strconv.Atoi(os.Getenv("REPEAT_SENT_CODE_TTL"))
	if err != nil {
		repeatSentCodeTtl = config.REPEAT_SENT_CODE_TTL
	}

	appConfig := config.AppConfig{
		CodeTtl:           codeTtl,
		RepeatSentCodeTtl: repeatSentCodeTtl,
	}

	return &appConfig
}
