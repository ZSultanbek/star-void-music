package repository

import (
	"context"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserLibraryRepository struct {
	queries *db.Queries
}

func NewUserLibraryRepository(pool *pgxpool.Pool) *UserLibraryRepository {
	return &UserLibraryRepository{queries: db.New(pool)}
}

func (r *UserLibraryRepository) AddSongToUserLibrary(ctx context.Context, params db.AddSongToUserLibraryParams) (db.UserLibrary, error) {
	return r.queries.AddSongToUserLibrary(ctx, params)
}

func (r *UserLibraryRepository) ListUserLibrarySongs(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.ListUserLibrarySongsRow, error) {
	return r.queries.ListUserLibrarySongs(ctx, db.ListUserLibrarySongsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *UserLibraryRepository) RemoveSongFromUserLibrary(ctx context.Context, params db.RemoveSongFromUserLibraryParams) error {
	return r.queries.RemoveSongFromUserLibrary(ctx, params)
}
