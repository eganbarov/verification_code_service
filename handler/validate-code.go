package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

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
