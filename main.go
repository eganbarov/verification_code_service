package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /", &homeHandler{})
	mux.Handle("POST /send-code", &SendCodeHandler{})
	mux.Handle("POST /validate-code", &ValidateCodeHandler{})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// Home page.
type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("phone")
	fmt.Println(id)
	w.Write([]byte("This is a home endpoint!"))
}

// Send code API.
type SendCodePostData struct {
	Phone string `json:"phone"`
}

type SentCodeResponse struct {
	IsSent     bool   `json:"isSent"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type SendCodeHandler struct{}

func (s *SendCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sendCodePost SendCodePostData
	err := json.NewDecoder(r.Body).Decode(&sendCodePost)
	if err != nil {
		http.Error(w, "Error decoding JSON send code data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	fmt.Println(sendCodePost.Phone)

	if sendCodePost.Phone == "" {
		http.Error(w, "Empty phone for send code", http.StatusBadRequest)
		return
	}

	code := rand.IntN(900000) + 100000

	//need to write it into redis and write to logs.
	fileName := sendCodePost.Phone + ".txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	codeString := strconv.Itoa(code)
	_, err = file.WriteString(codeString + "\n")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(code)

	isSentCode := true
	responseData := SentCodeResponse{
		IsSent:     isSentCode,
		StatusCode: http.StatusCreated,
		Error:      "",
	}
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Validate code API.
type ValidateCodePostData struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
}

type ValidateCodeResponse struct {
	IsValid    bool   `json:"isValid"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type ValidateCodeHandler struct{}

func (v *ValidateCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var validateCodePost ValidateCodePostData
	err := json.NewDecoder(r.Body).Decode(&validateCodePost)
	if err != nil {
		http.Error(w, "Error decoding JSON validate code data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	fmt.Println(validateCodePost.Phone)
	fmt.Println(validateCodePost.Code)

	if validateCodePost.Phone == "" {
		http.Error(w, "Error phone for validate code", http.StatusInternalServerError)
		return
	}

	if validateCodePost.Code == "" {
		http.Error(w, "Error code for validate code", http.StatusInternalServerError)
		return
	}

	//need to read it from redis.
	fileName := validateCodePost.Phone + ".txt"
	data, err := os.ReadFile(fileName)
	if err != nil {
		http.Error(w, "Error during reading file", http.StatusInternalServerError)
		return
	}

	code := string(data)
	if code != validateCodePost.Code {
		fmt.Println(code)
		fmt.Println(validateCodePost.Code)
		http.Error(w, "The code is not valid", http.StatusInternalServerError)
		return
	}

	responseData := ValidateCodeResponse{
		IsValid:    true,
		StatusCode: http.StatusOK,
		Error:      "",
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
