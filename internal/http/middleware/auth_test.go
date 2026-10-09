package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Max-shiri-90/OrderPulse/internal/http/middleware"
	"github.com/Max-shiri-90/OrderPulse/internal/user"
)

func TestRequireAuthMissingHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := middleware.RequireAuth("test-secret", next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRequireAuthInvalidToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := middleware.RequireAuth("test-secret", next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer invalid-token",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRequireAuthValidToken(t *testing.T) {
	token, err := user.GenerateToken(17, "test-secret")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		userID, ok := user.UserIDFromContext(r.Context())

		if !ok {
			t.Fatal("expected user ID in context")
		}

		if userID != 17 {
			t.Fatalf(
				"expected user ID 17, got %d",
				userID,
			)
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth("test-secret", next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestRequireAuthWrongSecret(t *testing.T) {
	token, err := user.GenerateToken(17, "correct-secret")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := middleware.RequireAuth("wrong-secret", next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}
