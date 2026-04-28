// User repository

package repository

import (
	"context"
	"time"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	queries *db.Queries
}

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: db.New(pool),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash, role string) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	})
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *UserRepository) GetUsers(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}