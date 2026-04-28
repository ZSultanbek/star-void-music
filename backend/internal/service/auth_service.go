package service

import (
	"context"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	users    *UserService
	secret   []byte
	tokenTTL time.Duration
}

type AuthResult struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewAuthService(users *UserService, secret []byte, tokenTTL time.Duration) (*AuthService, error) {
	if users == nil {
		return nil, errors.New("user service is required")
	}
	if len(secret) < 32 {
		return nil, errors.New("jwt secret must be at least 32 bytes")
	}
	if tokenTTL <= 0 {
		tokenTTL = time.Hour
	}
	return &AuthService{
		users:    users,
		secret:   secret,
		tokenTTL: tokenTTL,
	}, nil
}

func NewAuthServiceFromEnv(users *UserService, tokenTTL time.Duration) (*AuthService, error) {
	secret, err := ParseJWTSecretFromEnv()
	if err != nil {
		return nil, err
	}
	return NewAuthService(users, secret, tokenTTL)
}

func ParseJWTSecretFromEnv() ([]byte, error) {
	raw := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if raw == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	decoded, err := hex.DecodeString(raw)
	if err != nil {
		return nil, errors.New("JWT_SECRET must be a valid hex string")
	}
	if len(decoded) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 bytes")
	}
	return decoded, nil
}

func (s *AuthService) Register(ctx context.Context, email, password string) (AuthResult, error) {
	user, err := s.users.CreateUser(ctx, CreateUserInput{
		Email:        strings.TrimSpace(email),
		PasswordHash: password,
		Role:         "user",
	})
	if err != nil {
		return AuthResult{}, err
	}

	token, expiresIn, err := s.issueAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (AuthResult, error) {
	user, err := s.users.AuthenticateUser(ctx, email, password)
	if err != nil {
		return AuthResult{}, err
	}

	token, expiresIn, err := s.issueAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *AuthService) issueAccessToken(userID, role string) (string, int64, error) {
	now := time.Now().UTC()
	exp := now.Add(s.tokenTTL)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     exp.Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}
	return signed, int64(s.tokenTTL.Seconds()), nil
}
