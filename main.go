package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/eganbarov/verification_code_service/config"
	"github.com/eganbarov/verification_code_service/generator"
	"github.com/eganbarov/verification_code_service/handler"
	"github.com/eganbarov/verification_code_service/lock"
	"github.com/eganbarov/verification_code_service/middleware"
	"github.com/eganbarov/verification_code_service/repository"
	"github.com/eganbarov/verification_code_service/sender"
	"github.com/eganbarov/verification_code_service/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	startServer()
}

func startServer() {
	appConfig := loadAppConfig()
	db := initStorage(&appConfig.StorageConfig)
	codeRepository := repository.CodeRepository{
		Redis:     db,
		AppConfig: appConfig,
	}
	locker := lock.RedisLocker{
		Redis:     db,
		AppConfig: appConfig,
	}
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	wrappedMux := middleware.LoggingMiddleware(mux, logger)

	listenPort := os.Getenv("LISTEN_PORT")
	fmt.Println("Server starting on :" + listenPort)

	if err := http.ListenAndServe(":"+listenPort, wrappedMux); err != nil {
		log.Fatal(err)
	}
}

func initStorage(cnf *config.StorageConfig) *redis.Client {
	db, err := storage.NewClient(context.Background(), cnf)
	if err != nil {
		panic(err)
	}

	return db
}

func loadAppConfig() *config.AppConfig {
	appConfig := &config.AppConfig{}
	cnf, err := appConfig.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return cnf
}
