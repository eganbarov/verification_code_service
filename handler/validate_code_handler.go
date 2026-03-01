package handler

import (
	"encoding/json"
	"net/http"

	"github.com/eganbarov/verification_code_service/repository"
)

type validateCodePostData struct {
	Code   string `json:"code"`
	Phone  string `json:"phone"`
	Action string `json:"action"`
}

type validateCodeResponse struct {
	IsValid    bool   `json:"isValid"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type ValidateCodeHandler struct {
	CodeRepository repository.CodeRepo
}

func (v *ValidateCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var validateCodePost validateCodePostData
	err := json.NewDecoder(r.Body).Decode(&validateCodePost)
	if err != nil {
		renderErrorValidateCode(w, "Error decoding JSON validate code data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if validateCodePost.Phone == "" {
		renderErrorValidateCode(w, "Error phone for validate code", http.StatusBadRequest)
		return
	}

	if validateCodePost.Code == "" {
		renderErrorValidateCode(w, "Error code for validate code", http.StatusBadRequest)
		return
	}

	if validateCodePost.Action == "" {
		renderErrorValidateCode(w, "Error action for validate code", http.StatusBadRequest)
		return
	}

	value, err := v.CodeRepository.GetCode(validateCodePost.Phone, validateCodePost.Action)
	if err != nil {
		renderErrorValidateCode(w, "The code does not exist", http.StatusUnprocessableEntity)
		return
	}

	if value != validateCodePost.Code {
		renderErrorValidateCode(w, "The code is invalid", http.StatusUnprocessableEntity)
		return
	}

	if err := v.CodeRepository.DeleteCode(validateCodePost.Phone, validateCodePost.Action); err != nil {
		renderErrorValidateCode(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderSuccessValidateCode(w)
}

func renderSuccessValidateCode(w http.ResponseWriter) {
	responseData := validateCodeResponse{
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

func renderErrorValidateCode(w http.ResponseWriter, errMsg string, httpStatusCode int) {
	responseData := validateCodeResponse{
		IsValid:    false,
		StatusCode: httpStatusCode,
		Error:      errMsg,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(jsonData)
}
