package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type ValidateCodePostData struct {
	Code   string `json:"code"`
	Phone  string `json:"phone"`
	Action string `json:"action"`
}

type ValidateCodeResponse struct {
	IsValid    bool   `json:"isValid"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type ValidateCodeHandler struct {
	Redis *redis.Client
}

func (v *ValidateCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var validateCodePost ValidateCodePostData
	err := json.NewDecoder(r.Body).Decode(&validateCodePost)
	if err != nil {
		http.Error(w, "Error decoding JSON validate code data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if validateCodePost.Phone == "" {
		http.Error(w, "Error phone for validate code", http.StatusInternalServerError)
		return
	}

	if validateCodePost.Code == "" {
		http.Error(w, "Error code for validate code", http.StatusInternalServerError)
		return
	}

	if validateCodePost.Action == "" {
		http.Error(w, "Error action for validate code", http.StatusInternalServerError)
		return
	}

	codeKey := validateCodePost.Phone + "_" + validateCodePost.Action
	value, err := v.Redis.Get(context.Background(), codeKey).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if value != validateCodePost.Code {
		http.Error(w, "The code is invalid", http.StatusInternalServerError)
		return
	}

	if err := v.Redis.Del(context.Background(), codeKey).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
