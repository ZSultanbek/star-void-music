package repository

import (
	"context"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumRepository struct {
	queries *db.Queries
}

func NewAlbumRepository(pool *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{queries: db.New(pool)}
}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, params db.CreateAlbumParams) (db.Album, error) {
	return r.queries.CreateAlbum(ctx, params)
}

func (r *AlbumRepository) GetAlbumByID(ctx context.Context, id uuid.UUID) (db.Album, error) {
	return r.queries.GetAlbumByID(ctx, id)
}

func (r *AlbumRepository) ListAlbums(ctx context.Context, limit, offset int32) ([]db.Album, error) {
	return r.queries.ListAlbums(ctx, db.ListAlbumsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *AlbumRepository) ListAlbumsByArtistID(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Album, error) {
	return r.queries.ListAlbumsByArtistID(ctx, db.ListAlbumsByArtistIDParams{
		ArtistID: artistID,
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *AlbumRepository) UpdateAlbum(ctx context.Context, params db.UpdateAlbumParams) (db.Album, error) {
	return r.queries.UpdateAlbum(ctx, params)
}

func (r *AlbumRepository) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAlbum(ctx, id)
}

func (r *AlbumRepository) ListAlbumSongs(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]db.Song, error) {
	return r.queries.ListAlbumSongs(ctx, db.ListAlbumSongsParams{
		AlbumID: albumID,
		Limit:   limit,
		Offset:  offset,
	})
}
