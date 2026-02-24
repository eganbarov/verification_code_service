package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	mux.Handle("POST /send-code", &SendCodeHandler{})
	mux.Handle("POST /validate-code", &ValidateCodeHandler{})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
