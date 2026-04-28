package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var artistSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type ArtistRepository interface {
	CreateArtist(ctx context.Context, params db.CreateArtistParams) (db.Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (db.Artist, error)
	ListArtists(ctx context.Context, limit, offset int32) ([]db.Artist, error)
	UpdateArtist(ctx context.Context, params db.UpdateArtistParams) (db.Artist, error)
	DeleteArtist(ctx context.Context, id uuid.UUID) error
	ListArtistAlbums(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Album, error)
	ListArtistSongs(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Song, error)
}

type ArtistService struct {
	repo ArtistRepository
}

func NewArtistService(repo ArtistRepository) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) CreateArtist(ctx context.Context, params db.CreateArtistParams) (db.Artist, error) {
	params.Name = strings.TrimSpace(params.Name)
	params.Slug = strings.TrimSpace(strings.ToLower(params.Slug))
	if params.Name == "" || !artistSlugPattern.MatchString(params.Slug) {
		return db.Artist{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.CreateArtist(dbCtx, params)
}

func (s *ArtistService) GetArtistByID(ctx context.Context, id uuid.UUID) (db.Artist, error) {
	if id == uuid.Nil {
		return db.Artist{}, ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	artist, err := s.repo.GetArtistByID(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return db.Artist{}, ErrNotFound
	}
	return artist, err
}

func (s *ArtistService) ListArtists(ctx context.Context, limit, offset int32) ([]db.Artist, error) {
	limit, offset = normalizePagination(limit, offset)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListArtists(dbCtx, limit, offset)
}

func (s *ArtistService) UpdateArtist(ctx context.Context, params db.UpdateArtistParams) (db.Artist, error) {
	params.Name = strings.TrimSpace(params.Name)
	params.Slug = strings.TrimSpace(strings.ToLower(params.Slug))
	if params.ID == uuid.Nil || params.Name == "" || !artistSlugPattern.MatchString(params.Slug) {
		return db.Artist{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.UpdateArtist(dbCtx, params)
}

func (s *ArtistService) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.DeleteArtist(dbCtx, id)
}

func (s *ArtistService) ListArtistAlbums(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Album, error) {
	if artistID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListArtistAlbums(dbCtx, artistID, limit, offset)
}

func (s *ArtistService) ListArtistSongs(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Song, error) {
	if artistID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListArtistSongs(dbCtx, artistID, limit, offset)
}
