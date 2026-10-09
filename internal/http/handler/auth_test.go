package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
	"github.com/Max-shiri-90/OrderPulse/internal/user/memory"
)

func TestRegister(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	requestBody := `{
		"email": "test@example.com",
		"password": "my-secure-password"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Register(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	var response userResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != 1 {
		t.Fatalf("expected ID 1, got %d", response.ID)
	}

	if response.Email != "test@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"test@example.com",
			response.Email,
		)
	}

	if strings.Contains(
		recorder.Body.String(),
		"my-secure-password",
	) {
		t.Fatal("response must not contain the password")
	}
}

func TestRegisterInvalidRequestBody(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		strings.NewReader(`invalid-json`),
	)

	recorder := httptest.NewRecorder()

	handler.Register(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestRegisterInvalidEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	requestBody := `{
		"email": "",
		"password": "my-secure-password"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Register(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	requestBody := `{
		"email": "test@example.com",
		"password": ""
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Register(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	requestBody := `{
		"email": "test@example.com",
		"password": "another-password"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Register(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusConflict,
			recorder.Code,
		)
	}
}

func TestLogin(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	_, err := service.Register(
		context.Background(),
		"login@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	requestBody := `{
		"email": "login@example.com",
		"password": "my-secure-password"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Login(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var response loginResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Token == "" {
		t.Fatal("expected a non-empty JWT")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	_, err := service.Register(
		context.Background(),
		"login@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	requestBody := `{
		"email": "login@example.com",
		"password": "wrong-password"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Login(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusUnauthorized,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	requestBody := `{
		"email": "unknown@example.com",
		"password": "my-secure-password"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.Login(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusUnauthorized,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestLoginInvalidRequestBody(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		strings.NewReader(`invalid-json`),
	)

	recorder := httptest.NewRecorder()

	handler.Login(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestGetMeUnauthenticated(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/me",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetMe(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusUnauthorized,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestGetMeAuthenticated(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)
	handler := NewAuthHandler(service, "test-secret")

	createdUser, err := service.Register(
		context.Background(),
		"me@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/me",
		nil,
	)

	request = request.WithContext(
		user.ContextWithUserID(
			request.Context(),
			createdUser.ID,
		),
	)

	recorder := httptest.NewRecorder()

	handler.GetMe(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var response userResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != createdUser.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			createdUser.ID,
			response.ID,
		)
	}

	if response.Email != "me@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"me@example.com",
			response.Email,
		)
	}

	if strings.Contains(recorder.Body.String(), "my-secure-password") {
		t.Fatal("response must not contain the password")
	}
}

