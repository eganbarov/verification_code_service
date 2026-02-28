package handler

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type SendCodePostData struct {
	Phone  string `json:"phone"`
	Action string `json:"action"`
}

type SentCodeResponse struct {
	IsSent     bool   `json:"isSent"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type SendCodeHandler struct {
	Redis *redis.Client
}

func (s *SendCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sendCodePost SendCodePostData
	err := json.NewDecoder(r.Body).Decode(&sendCodePost)
	if err != nil {
		http.Error(w, "Error decoding JSON send code data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if sendCodePost.Phone == "" {
		http.Error(w, "Empty phone for send code", http.StatusBadRequest)
		return
	}

	if sendCodePost.Action == "" {
		http.Error(w, "Empty action for send code", http.StatusBadRequest)
		return
	}

	codeKey := sendCodePost.Phone + "_" + sendCodePost.Action
	code := rand.IntN(900000) + 100000
	if err := s.Redis.Set(context.Background(), codeKey, code, 300*time.Second).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send code via SMS

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
