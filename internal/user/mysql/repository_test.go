package mysql

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
)

func testDatabase(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "orderpulse:orderpulsepassword@tcp(127.0.0.1:3307)/orderpulse?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	u := user.User{
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
	}

	created, err := repository.Create(
		context.Background(),
		u,
	)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if created.ID == 0 {
		t.Fatal("expected created user to have an ID")
	}

	if created.Email != u.Email {
		t.Fatalf("expected email %q, got %q", u.Email, created.Email)
	}

	if created.PasswordHash != u.PasswordHash {
		t.Fatalf(
			"expected password hash %q, got %q",
			u.PasswordHash,
			created.PasswordHash,
		)
	}
}

func TestGetByID(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	created, err := repository.Create(
		context.Background(),
		user.User{
			Email:        "id-test@example.com",
			PasswordHash: "hashed-password",
		},
	)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	found, err := repository.GetByID(
		context.Background(),
		created.ID,
	)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if found.ID != created.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			created.ID,
			found.ID,
		)
	}

	if found.Email != created.Email {
		t.Fatalf(
			"expected email %q, got %q",
			created.Email,
			found.Email,
		)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	_, err = repository.GetByID(
		context.Background(),
		999999,
	)

	if err != user.ErrNotFound {
		t.Fatalf(
			"expected ErrNotFound, got %v",
			err,
		)
	}
}

func TestGetByEmail(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	created, err := repository.Create(
		context.Background(),
		user.User{
			Email:        "email-test@example.com",
			PasswordHash: "hashed-password",
		},
	)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	found, err := repository.GetByEmail(
		context.Background(),
		created.Email,
	)
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	if found.ID != created.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			created.ID,
			found.ID,
		)
	}

	if found.Email != created.Email {
		t.Fatalf(
			"expected email %q, got %q",
			created.Email,
			found.Email,
		)
	}

	if found.PasswordHash != created.PasswordHash {
		t.Fatalf(
			"expected password hash %q, got %q",
			created.PasswordHash,
			found.PasswordHash,
		)
	}
}

func TestGetByEmailNotFound(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	_, err = repository.GetByEmail(
		context.Background(),
		"missing@example.com",
	)

	if err != user.ErrNotFound {
		t.Fatalf(
			"expected ErrNotFound, got %v",
			err,
		)
	}
}
