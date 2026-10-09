package user

import "testing"

func TestHashPassword(t *testing.T) {
	password := "my-secure-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected password hash to be non-empty")
	}

	if hash == password {
		t.Fatal("password must not be stored as plain text")
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Fatalf("expected password to match hash: %v", err)
	}
}

func TestCheckPasswordWrongPassword(t *testing.T) {
	password := "my-secure-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = CheckPassword("wrong-password", hash)

	if err == nil {
		t.Fatal("expected wrong password to fail")
	}
}

func TestHashPasswordProducesDifferentHashes(t *testing.T) {
	password := "my-secure-password"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to generate first hash: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to generate second hash: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("expected different hashes for the same password")
	}
}
