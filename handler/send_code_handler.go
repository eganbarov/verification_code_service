package handler

import (
	"encoding/json"
	"net/http"

	"github.com/eganbarov/verification_code_service/generator"
	"github.com/eganbarov/verification_code_service/lock"
	"github.com/eganbarov/verification_code_service/repository"
	"github.com/eganbarov/verification_code_service/sender"
)

type sendCodePostData struct {
	Phone  string `json:"phone"`
	Action string `json:"action"`
}

type sentCodeResponse struct {
	IsSent     bool   `json:"isSent"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type SendCodeHandler struct {
	CodeRepository repository.CodeRepo
	Locker         lock.Locker
	CodeGenerator  generator.CodeGen
	CodeSender     sender.CodeSender
}

func (s *SendCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sendCodePost sendCodePostData
	err := json.NewDecoder(r.Body).Decode(&sendCodePost)
	if err != nil {
		renderErrorSentCode(w, "Error decoding JSON send code data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if sendCodePost.Phone == "" {
		renderErrorSentCode(w, "Empty phone for send code", http.StatusBadRequest)
		return
	}

	if sendCodePost.Action == "" {
		renderErrorSentCode(w, "Empty action for send code", http.StatusBadRequest)
		return
	}

	isLocked := s.Locker.IsLocked(sendCodePost.Phone, sendCodePost.Action)
	if isLocked == true {
		renderSuccessSentCode(w)
		return
	}

	code := s.CodeGenerator.GenerateCode()
	if err := s.CodeRepository.StoreCode(sendCodePost.Phone, sendCodePost.Action, code); err != nil {
		renderErrorSentCode(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.CodeSender.SendCode(code); err != nil {
		renderErrorSentCode(w, "Error during sending a code", http.StatusInternalServerError)
		return
	}

	if err := s.Locker.Lock(sendCodePost.Phone, sendCodePost.Action); err != nil {
		renderErrorSentCode(w, "Error during lock sending", http.StatusInternalServerError)
		return
	}

	renderSuccessSentCode(w)
}

func renderSuccessSentCode(w http.ResponseWriter) {
	responseData := sentCodeResponse{
		IsSent:     true,
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

func renderErrorSentCode(w http.ResponseWriter, errMsg string, httpStatusCode int) {
	responseData := sentCodeResponse{
		IsSent:     false,
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
