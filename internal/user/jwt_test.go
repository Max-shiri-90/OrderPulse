package user

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret"
	userID := int64(42)

	tokenString, err := GenerateToken(userID, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected token to be non-empty")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Fatal("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected MapClaims")
	}

	claimUserID, ok := claims["user_id"].(float64)
	if !ok {
		t.Fatal("expected user_id claim to be a number")
	}

	if int64(claimUserID) != userID {
		t.Fatalf(
			"expected user_id %d, got %d",
			userID,
			int64(claimUserID),
		)
	}
}

func TestGenerateTokenHasExpiration(t *testing.T) {
	secret := "test-secret"

	tokenString, err := GenerateToken(42, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected MapClaims")
	}

	expiration, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("expected exp claim")
	}

	expirationTime := time.Unix(int64(expiration), 0)

	if !expirationTime.After(time.Now()) {
		t.Fatal("expected token expiration to be in the future")
	}
}

func TestGenerateTokenRejectsWrongSecret(t *testing.T) {
	tokenString, err := GenerateToken(42, "correct-secret")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte("wrong-secret"), nil
		},
	)

	if err == nil {
		t.Fatal("expected token validation to fail with wrong secret")
	}

	if token != nil && token.Valid {
		t.Fatal("expected token to be invalid")
	}
}
