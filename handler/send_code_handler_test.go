package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/eganbarov/verification_code_service/config"
	"github.com/eganbarov/verification_code_service/generator"
	"github.com/eganbarov/verification_code_service/handler"
	"github.com/eganbarov/verification_code_service/lock"
	"github.com/eganbarov/verification_code_service/repository"
	"github.com/eganbarov/verification_code_service/sender"
	"github.com/redis/go-redis/v9"
)

func TestSucces_SendCodeServeHTTP(t *testing.T) {
	mr := miniredis.RunT(t)
	storage := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer storage.Close()

	phone := "+79160010203"
	action := "auth"

	postData := map[string]string{
		"phone":  phone,
		"action": action,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		t.Fatal(err)
	}

	reqBody := bytes.NewReader(jsonData)
	req := httptest.NewRequest("POST", "/send-code", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	appCnf := &config.AppConfig{}
	locker := &lock.RedisLocker{Redis: storage, AppConfig: appCnf}
	repo := &repository.CodeRepository{Redis: storage, AppConfig: appCnf}

	handler := handler.SendCodeHandler{
		CodeRepository: repo,
		Locker:         locker,
		CodeGenerator:  &generator.CodeGenerator{},
		CodeSender:     &sender.SmsSender{},
	}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, rec.Code)
	}

	expectedBody := `{"isSent":true,"statusCode":201,"error":""}`
	gotBody := rec.Body.String()
	if gotBody != expectedBody {
		t.Errorf("Expected body: %v, got: , %v", expectedBody, gotBody)
	}
}

func TestValidation_SendCodeServeHTTP(t *testing.T) {
	tests := map[string]struct {
		phone   string
		action  string
		recCode int
		recBody string
	}{
		"empty params": {
			phone:   "",
			action:  "",
			recCode: http.StatusBadRequest,
		},
		"empty phone": {
			phone:   "",
			action:  "auth",
			recCode: http.StatusBadRequest,
		},
		"empty action": {
			phone:   "+79160010203",
			action:  "",
			recCode: http.StatusBadRequest,
		},
	}

	mr := miniredis.RunT(t)
	storage := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer storage.Close()

	appCnf := &config.AppConfig{}
	locker := &lock.RedisLocker{Redis: storage, AppConfig: appCnf}
	repo := &repository.CodeRepository{Redis: storage, AppConfig: appCnf}

	handler := handler.SendCodeHandler{
		CodeRepository: repo,
		Locker:         locker,
		CodeGenerator:  &generator.CodeGenerator{},
		CodeSender:     &sender.SmsSender{},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			postData := map[string]string{
				"phone":  test.phone,
				"action": test.action,
			}
			jsonData, err := json.Marshal(postData)
			if err != nil {
				t.Fatal(err)
			}

			reqBody := bytes.NewReader(jsonData)
			req := httptest.NewRequest("POST", "/send-code", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != test.recCode {
				t.Errorf("Expected %d, got %d", test.recCode, rec.Code)
			}
		})
	}
}
