package middleware

import (
	"net/http"
	"strings"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(
	jwtSecret string,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"missing authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			http.Error(
				w,
				"invalid authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			jwt.MapClaims{},
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(
				w,
				"invalid or expired token",
				http.StatusUnauthorized,
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(
				w,
				"invalid token claims",
				http.StatusUnauthorized,
			)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok || userID <= 0 {
			http.Error(
				w,
				"invalid token claims",
				http.StatusUnauthorized,
			)
			return
		}

		ctx := user.ContextWithUserID(
			r.Context(),
			int64(userID),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
