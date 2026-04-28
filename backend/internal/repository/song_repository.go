package repository

import (
	"context"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepository struct {
	queries *db.Queries
}

func NewSongRepository(pool *pgxpool.Pool) *SongRepository {
	return &SongRepository{queries: db.New(pool)}
}

func (r *SongRepository) CreateSong(ctx context.Context, params db.CreateSongParams) (db.Song, error) {
	return r.queries.CreateSong(ctx, params)
}

func (r *SongRepository) GetSongByID(ctx context.Context, id uuid.UUID) (db.Song, error) {
	return r.queries.GetSongByID(ctx, id)
}

func (r *SongRepository) ListSongs(ctx context.Context, limit, offset int32) ([]db.Song, error) {
	return r.queries.ListSongs(ctx, db.ListSongsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *SongRepository) ListSongsByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]db.Song, error) {
	return r.queries.ListSongsByAlbumID(ctx, db.ListSongsByAlbumIDParams{
		AlbumID: albumID,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *SongRepository) UpdateSong(ctx context.Context, params db.UpdateSongParams) (db.Song, error) {
	return r.queries.UpdateSong(ctx, params)
}

func (r *SongRepository) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSong(ctx, id)
}
