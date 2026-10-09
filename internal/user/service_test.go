package user_test

import (
	"context"
	"testing"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
	"github.com/Max-shiri-90/OrderPulse/internal/user/memory"
)

func TestRegister(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	created, err := service.Register(
		context.Background(),
		"test@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}

	if created.Email != "test@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"test@example.com",
			created.Email,
		)
	}

	if created.PasswordHash == "" {
		t.Fatal("expected password hash to be set")
	}

	if created.PasswordHash == "my-secure-password" {
		t.Fatal("password must not be stored as plain text")
	}

	if err := user.CheckPassword(
		"my-secure-password",
		created.PasswordHash,
	); err != nil {
		t.Fatalf("expected password to match hash: %v", err)
	}
}

func TestRegisterTrimsEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	created, err := service.Register(
		context.Background(),
		"  test@example.com  ",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.Email != "test@example.com" {
		t.Fatalf(
			"expected trimmed email %q, got %q",
			"test@example.com",
			created.Email,
		)
	}
}

func TestRegisterInvalidEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	_, err := service.Register(
		context.Background(),
		"",
		"my-secure-password",
	)

	if err != user.ErrInvalidEmail {
		t.Fatalf(
			"expected ErrInvalidEmail, got %v",
			err,
		)
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"",
	)

	if err != user.ErrInvalidPassword {
		t.Fatalf(
			"expected ErrInvalidPassword, got %v",
			err,
		)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	_, err = service.Register(
		context.Background(),
		"test@example.com",
		"another-password",
	)

	if err != user.ErrEmailExists {
		t.Fatalf(
			"expected ErrEmailExists, got %v",
			err,
		)
	}
}

func TestLogin(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	token, err := service.Login(
		context.Background(),
		"test@example.com",
		"my-secure-password",
		"test-secret",
	)
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}

	if token == "" {
		t.Fatal("expected JWT token to be non-empty")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"my-secure-password",
	)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	_, err = service.Login(
		context.Background(),
		"test@example.com",
		"wrong-password",
		"test-secret",
	)

	if err != user.ErrInvalidCredentials {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	repository := memory.NewRepository()
	service := user.NewService(repository)

	_, err := service.Login(
		context.Background(),
		"unknown@example.com",
		"my-secure-password",
		"test-secret",
	)

	if err != user.ErrInvalidCredentials {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}
