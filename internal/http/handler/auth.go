package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
)

type AuthHandler struct {
	service   *user.Service
	jwtSecret string
}

func NewAuthHandler(service *user.Service, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func newUserResponse(u user.User) userResponse {
	return userResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	u, err := h.service.Register(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, user.ErrInvalidEmail) ||
			errors.Is(err, user.ErrInvalidPassword) {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		if errors.Is(err, user.ErrEmailExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)
			return
		}

		http.Error(
			w,
			"failed to register user",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := newUserResponse(u)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	token, err := h.service.Login(
		r.Context(),
		req.Email,
		req.Password,
		h.jwtSecret,
	)

	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			http.Error(
				w,
				err.Error(),
				http.StatusUnauthorized,
			)
			return
		}

		http.Error(
			w,
			"failed to login",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := loginResponse{
		Token: token,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := user.UserIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	u, err := h.service.GetByID(
		r.Context(),
		userID,
	)

	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			http.Error(
				w,
				"user not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"failed to get user",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := newUserResponse(u)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
