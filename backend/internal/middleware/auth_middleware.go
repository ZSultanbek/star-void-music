package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
	"bytes"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

type authUserContextKey struct{}

type AuthenticatedUser struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

func NewJWTAuthMiddlewareFromEnv() (gin.HandlerFunc, error) {
	secret := os.Getenv("JWT_SECRET")
	secretbytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, errors.New("JWT_SECRET must be a valid hex string")
	}
	if len(secretbytes) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 bytes")
	}
	secretbytes = bytes.TrimSpace(secretbytes)
	secret = string(secretbytes)
	return JWTAuthMiddleware([]byte(secret)), nil
}

func JWTAuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c.GetHeader(authorizationHeader))
		if err != nil {
			abortUnauthorized(c)
			return
		}

		user, err := ParseAndValidateAccessToken(tokenString, secret, time.Now())
		if err != nil {
			abortUnauthorized(c)
			return
		}

		// type-safe context storage
		ctx := context.WithValue(c.Request.Context(), authUserContextKey{}, user)
		c.Request = c.Request.WithContext(ctx)

		// compatibility with existing handlers
		c.Set("user_id", user.UserID)
		c.Set("role", user.Role)

		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	requiredRole = strings.ToLower(strings.TrimSpace(requiredRole))

	return func(c *gin.Context) {
		user, ok := GetUserFromContext(c.Request.Context())
		if !ok {
			abortUnauthorized(c)
			return
		}
		if strings.ToLower(strings.TrimSpace(user.Role)) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"data":  nil,
				"error": "forbidden",
			})
			return
		}
		c.Next()
	}
}

func GetUserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	v := ctx.Value(authUserContextKey{})
	user, ok := v.(AuthenticatedUser)
	return user, ok
}

func ParseAndValidateAccessToken(tokenString string, secret []byte, now time.Time) (AuthenticatedUser, error) {
	if strings.TrimSpace(tokenString) == "" || len(secret) == 0 {
		return AuthenticatedUser{}, ErrUnauthorized
	}

	claims := &AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrUnauthorized
			}
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(30*time.Second),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	if err != nil || token == nil || !token.Valid {
		return AuthenticatedUser{}, ErrUnauthorized
	}

	userID, err := uuid.Parse(strings.TrimSpace(claims.UserID))
	if err != nil || userID == uuid.Nil {
		return AuthenticatedUser{}, ErrUnauthorized
	}

	role := strings.ToLower(strings.TrimSpace(claims.Role))
	if role == "" {
		return AuthenticatedUser{}, ErrUnauthorized
	}

	return AuthenticatedUser{
		UserID: userID,
		Role:   role,
	}, nil
}

func extractBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", ErrUnauthorized
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		return "", ErrUnauthorized
	}

	return token, nil
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"data":  nil,
		"error": "unauthorized",
	})
}
