package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eganbarov/verification_code_service/handler"
)

func main() {
	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	mux.Handle("POST /send-code", &handler.SendCodeHandler{})
	mux.Handle("POST /validate-code", &handler.ValidateCodeHandler{})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
