package handler

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

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
