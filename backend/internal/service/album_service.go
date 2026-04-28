package service

import (
	"context"
	"errors"
	"strings"
	"time"

	db "star-void-music/backend/db/sqlc"
	"star-void-music/backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, params db.CreateAlbumParams) (db.Album, error)
	GetAlbumByID(ctx context.Context, id uuid.UUID) (db.Album, error)
	ListAlbums(ctx context.Context, limit, offset int32) ([]db.Album, error)
	ListAlbumsByArtistID(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Album, error)
	UpdateAlbum(ctx context.Context, params db.UpdateAlbumParams) (db.Album, error)
	DeleteAlbum(ctx context.Context, id uuid.UUID) error
	ListAlbumSongs(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]db.Song, error)
}

type AlbumService struct {
	repo AlbumRepository
}

func NewAlbumService(repo AlbumRepository) *AlbumService {
	return &AlbumService{repo: repo}
}

type CreateAlbumInput struct {
	Title         string
	ArtistID      uuid.UUID
	CoverImageURL string
	ReleaseDate   *time.Time
}

type UpdateAlbumInput struct {
	ID            uuid.UUID
	Title         string
	ArtistID      uuid.UUID
	CoverImageURL string
	ReleaseDate   *time.Time
}

func toPGDate(v *time.Time) pgtype.Date {
	if v == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *v, Valid: true}
}

func mapAlbumRow(r db.Album) models.Album {
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

func (s *AlbumService) CreateAlbum(ctx context.Context, in CreateAlbumInput) (models.Album, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.CoverImageURL = strings.TrimSpace(in.CoverImageURL)
	if in.Title == "" || in.ArtistID == uuid.Nil {
		return models.Album{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	row, err := s.repo.CreateAlbum(dbCtx, db.CreateAlbumParams{
		Title: in.Title, ArtistID: in.ArtistID, CoverImageUrl: in.CoverImageURL, ReleaseDate: toPGDate(in.ReleaseDate),
	})
	if err != nil {
		return models.Album{}, err
	}
	return mapAlbumRow(row), nil
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, id uuid.UUID) (models.Album, error) {
	if id == uuid.Nil {
		return models.Album{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row, err := s.repo.GetAlbumByID(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return models.Album{}, ErrNotFound
	}
	if err != nil {
		return models.Album{}, err
	}
	return mapAlbumRow(row), nil
}

func (s *AlbumService) ListAlbums(ctx context.Context, limit, offset int32) ([]db.Album, error) {
	limit, offset = normalizePagination(limit, offset)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListAlbums(dbCtx, limit, offset)
}

func (s *AlbumService) ListAlbumsByArtistID(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]db.Album, error) {
	if artistID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListAlbumsByArtistID(dbCtx, artistID, limit, offset)
}

func (s *AlbumService) UpdateAlbum(ctx context.Context, in UpdateAlbumInput) (models.Album, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.CoverImageURL = strings.TrimSpace(in.CoverImageURL)
	if in.ID == uuid.Nil || in.ArtistID == uuid.Nil || in.Title == "" {
		return models.Album{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	row, err := s.repo.UpdateAlbum(dbCtx, db.UpdateAlbumParams{
		ID: in.ID, Title: in.Title, ArtistID: in.ArtistID, CoverImageUrl: in.CoverImageURL, ReleaseDate: toPGDate(in.ReleaseDate),
	})
	if err != nil {
		return models.Album{}, err
	}
	return mapAlbumRow(row), nil
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.DeleteAlbum(dbCtx, id)
}

func (s *AlbumService) ListAlbumSongs(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]models.Song, error) {
	if albumID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	rows, err := s.repo.ListAlbumSongs(dbCtx, albumID, limit, offset)
	if err != nil {
		return nil, err
	}
	return mapSongs(rows), nil
}
