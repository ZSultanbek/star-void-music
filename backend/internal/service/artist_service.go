package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	db "star-void-music/backend/db/sqlc"
	"star-void-music/backend/internal/models"

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

type CreateArtistInput struct {
	Name string
	Slug string
}

type UpdateArtistInput struct {
	ID   uuid.UUID
	Name string
	Slug string
}

func (s *ArtistService) CreateArtist(ctx context.Context, in CreateArtistInput) (models.Artist, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Slug = strings.TrimSpace(strings.ToLower(in.Slug))
	if in.Name == "" || !artistSlugPattern.MatchString(in.Slug) {
		return models.Artist{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	row, err := s.repo.CreateArtist(dbCtx, db.CreateArtistParams{Name: in.Name, Slug: in.Slug})
	if err != nil {
		return models.Artist{}, err
	}
	return mapArtist(row), nil
}

func (s *ArtistService) GetArtistByID(ctx context.Context, id uuid.UUID) (models.Artist, error) {
	if id == uuid.Nil {
		return models.Artist{}, ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row, err := s.repo.GetArtistByID(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return models.Artist{}, ErrNotFound
	}
	if err != nil {
		return models.Artist{}, err
	}
	return mapArtist(row), nil
}

func (s *ArtistService) ListArtists(ctx context.Context, limit, offset int32) ([]models.Artist, error) {
	limit, offset = normalizePagination(limit, offset)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	rows, err := s.repo.ListArtists(dbCtx, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]models.Artist, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapArtist(r))
	}
	return out, nil
}

func (s *ArtistService) UpdateArtist(ctx context.Context, in UpdateArtistInput) (models.Artist, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Slug = strings.TrimSpace(strings.ToLower(in.Slug))
	if in.ID == uuid.Nil || in.Name == "" || !artistSlugPattern.MatchString(in.Slug) {
		return models.Artist{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	row, err := s.repo.UpdateArtist(dbCtx, db.UpdateArtistParams{ID: in.ID, Name: in.Name, Slug: in.Slug})
	if err != nil {
		return models.Artist{}, err
	}
	return mapArtist(row), nil
}

func (s *ArtistService) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.DeleteArtist(dbCtx, id)
}

func (s *ArtistService) ListArtistAlbums(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]models.Album, error) {
	if artistID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	rows, err := s.repo.ListArtistAlbums(dbCtx, artistID, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]models.Album, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapAlbum(r))
	}
	return out, nil
}

func mapArtist(r db.Artist) models.Artist {
	return models.Artist{ID: r.ID, Name: r.Name, Slug: r.Slug, CreatedAt: r.CreatedAt}
}

func mapAlbum(r db.Album) models.Album {
	var rd *time.Time
	if r.ReleaseDate.Valid {
		t := r.ReleaseDate.Time
		rd = &t
	}
	return models.Album{
		ID: r.ID, Title: r.Title, ArtistID: r.ArtistID,
		CoverImageURL: r.CoverImageUrl, ReleaseDate: rd, CreatedAt: r.CreatedAt,
	}
}
