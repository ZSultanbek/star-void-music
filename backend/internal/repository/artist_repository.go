package repository

import (
	"context"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistRepository struct {
	queries *db.Queries
}

func NewArtistRepository(pool *pgxpool.Pool) *ArtistRepository {
	return &ArtistRepository{queries: db.New(pool)}
}

func (r *ArtistRepository) CreateArtist(ctx context.Context, params db.CreateArtistParams) (db.Artist, error) {
	return r.queries.CreateArtist(ctx, params)
}

func (r *ArtistRepository) GetArtistByID(ctx context.Context, id uuid.UUID) (db.Artist, error) {
	return r.queries.GetArtistByID(ctx, id)
}

func (r *ArtistRepository) ListArtists(ctx context.Context, limit, offset int32) ([]db.Artist, error) {
	return r.queries.ListArtists(ctx, db.ListArtistsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *ArtistRepository) UpdateArtist(ctx context.Context, params db.UpdateArtistParams) (db.Artist, error) {
	return r.queries.UpdateArtist(ctx, params)
}

func (r *ArtistRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteArtist(ctx, id)
}

func (r *ArtistRepository) ListArtistAlbums(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Album, error) {
	return r.queries.ListArtistAlbums(ctx, db.ListArtistAlbumsParams{
		ArtistID: artistID,
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *ArtistRepository) ListArtistSongs(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Song, error) {
	return r.queries.ListArtistSongs(ctx, db.ListArtistSongsParams{
		ArtistID: artistID,
		Limit:    limit,
		Offset:   offset,
	})
}
