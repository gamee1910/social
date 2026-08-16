package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/pkg/utils"
	"github.com/gamee1910/social/pkg/jwt"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

var (
	errMissingToken       = errors.New("authorization token is required")
	errInvalidTokenFormat = errors.New("authorization header format must be 'Bearer <token>'")
)

func JWTAuth(jwtConfig config.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.Unauthorized(w, r, errMissingToken)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				utils.Unauthorized(w, r, errInvalidTokenFormat)
				return
			}

			claims, err := jwt.ValidateToken(parts[1], jwtConfig.SecretKey)
			if err != nil {
				utils.Unauthorized(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserClaims(r *http.Request) *jwt.Claims {
	claims, _ := r.Context().Value(UserClaimsKey).(*jwt.Claims)
	return claims
}
