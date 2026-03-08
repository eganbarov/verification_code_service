package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/eganbarov/verification_code_service/config"
	"github.com/eganbarov/verification_code_service/handler"
	"github.com/eganbarov/verification_code_service/lock"
	"github.com/eganbarov/verification_code_service/repository"
	"github.com/redis/go-redis/v9"
)

func TestSuccess_ValidateCodeServeHTTP(t *testing.T) {
	mr := miniredis.RunT(t)
	storage := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer storage.Close()

	phone := "+79160010203"
	action := "auth"
	code := "001002"
	postData := map[string]string{
		"code":   code,
		"phone":  phone,
		"action": action,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		t.Fatal(err)
	}

	reqBody := bytes.NewReader(jsonData)
	req := httptest.NewRequest("POST", "/validate-code", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	appCnf := &config.AppConfig{}
	locker := &lock.RedisLocker{Redis: storage, AppConfig: appCnf}
	repo := &repository.CodeRepository{Redis: storage, AppConfig: appCnf}
	repo.StoreCode(phone, action, code)

	handler := &handler.ValidateCodeHandler{CodeRepository: repo, Locker: locker}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := `{"isValid":true,"statusCode":200,"error":""}`
	gotBody := rec.Body.String()
	if gotBody != expectedBody {
		t.Errorf("Expected body: %v, got: , %v", expectedBody, gotBody)
	}
}

func TestValidation_ValidateCodeServeHTTP(t *testing.T) {
	tests := map[string]struct {
		phone   string
		action  string
		code    string
		recCode int
		recBody string
	}{
		"empty params": {
			phone:   "",
			action:  "",
			code:    "",
			recCode: http.StatusBadRequest,
		},
		"empty phone": {
			phone:   "",
			action:  "auth",
			code:    "001002",
			recCode: http.StatusBadRequest,
		},
		"empty action": {
			phone:   "+79160010203",
			action:  "",
			code:    "001002",
			recCode: http.StatusBadRequest,
		},
		"empty code": {
			phone:   "+79160010203",
			action:  "auth",
			code:    "",
			recCode: http.StatusBadRequest,
		},
		"code does not exist": {
			phone:   "+79160010203",
			action:  "auth",
			code:    "001002",
			recCode: http.StatusUnprocessableEntity,
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
	handler := &handler.ValidateCodeHandler{CodeRepository: repo, Locker: locker}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			postData := map[string]string{
				"code":   test.code,
				"phone":  test.phone,
				"action": test.action,
			}
			jsonData, err := json.Marshal(postData)
			if err != nil {
				t.Fatal(err)
			}

			reqBody := bytes.NewReader(jsonData)
			req := httptest.NewRequest("POST", "/validate-code", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != test.recCode {
				t.Errorf("Expected %d, got %d", test.recCode, rec.Code)
			}
		})
	}

}
