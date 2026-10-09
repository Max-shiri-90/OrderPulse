package user

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (User, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return User{}, ErrInvalidEmail
	}

	if password == "" {
		return User{}, ErrInvalidPassword
	}

	_, err := s.repository.GetByEmail(ctx, email)

	if err == nil {
		return User{}, ErrEmailExists
	}

	if !errors.Is(err, ErrNotFound) {
		return User{}, err
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	u := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	return s.repository.Create(ctx, u)
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
	secret string,
) (string, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return "", ErrInvalidCredentials
	}

	if password == "" {
		return "", ErrInvalidCredentials
	}

	u, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	if err := CheckPassword(password, u.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateToken(u.ID, secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (User, error) {
	if id <= 0 {
		return User{}, ErrNotFound
	}

	return s.repository.GetByID(ctx, id)
}
