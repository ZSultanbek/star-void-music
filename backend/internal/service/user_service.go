package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

const dbTimeout = 3 * time.Second

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash, role string) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUsers(ctx context.Context, limit, offset int32) ([]db.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	repo UserRepository
}

type CreateUserInput struct {
	Email        string
	PasswordHash string
	Role         string
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func normalizePagination(limit, offset int32) (int32, int32) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func sanitizeUser(u db.User) db.User {
	u.PasswordHash = ""
	return u
}

func sanitizeUsers(users []db.User) []db.User {
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users
}

func (s *UserService) CreateUser(ctx context.Context, in CreateUserInput) (db.User, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Role = strings.TrimSpace(strings.ToLower(in.Role))
	rawPassword := strings.TrimSpace(in.PasswordHash)

	if _, err := mail.ParseAddress(in.Email); err != nil {
		return db.User{}, ErrValidation
	}
	if len(rawPassword) < 8 {
		return db.User{}, ErrValidation
	}
	if in.Role == "" {
		in.Role = "user"
	}
	if in.Role != "user" && in.Role != "admin" {
		return db.User{}, ErrValidation
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := s.repo.CreateUser(dbCtx, in.Email, string(hashed), in.Role)
	if err != nil {
		return db.User{}, err
	}
	return sanitizeUser(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	if id == uuid.Nil {
		return db.User{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := s.repo.GetUserByID(dbCtx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return sanitizeUser(user), nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if _, err := mail.ParseAddress(email); err != nil {
		return db.User{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := s.repo.GetUserByEmail(dbCtx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return sanitizeUser(user), nil
}

func (s *UserService) GetUsers(ctx context.Context, limit, offset int32) ([]db.User, error) {
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	users, err := s.repo.GetUsers(dbCtx, limit, offset)
	if err != nil {
		return nil, err
	}
	return sanitizeUsers(users), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := s.repo.DeleteUser(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (db.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if _, err := mail.ParseAddress(email); err != nil || len(password) < 8 {
		return db.User{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := s.repo.GetUserByEmail(dbCtx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrUnauthorized
		}
		return db.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return db.User{}, ErrUnauthorized
	}

	return sanitizeUser(user), nil
}
