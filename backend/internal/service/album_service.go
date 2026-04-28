package service

import (
	"context"
	"errors"
	"strings"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *AlbumService) CreateAlbum(ctx context.Context, params db.CreateAlbumParams) (db.Album, error) {
	params.Title = strings.TrimSpace(params.Title)
	params.CoverImageUrl = strings.TrimSpace(params.CoverImageUrl)
	if params.Title == "" || params.ArtistID == uuid.Nil {
		return db.Album{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.CreateAlbum(dbCtx, params)
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, id uuid.UUID) (db.Album, error) {
	if id == uuid.Nil {
		return db.Album{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	album, err := s.repo.GetAlbumByID(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return db.Album{}, ErrNotFound
	}
	return album, err
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

func (s *AlbumService) UpdateAlbum(ctx context.Context, params db.UpdateAlbumParams) (db.Album, error) {
	params.Title = strings.TrimSpace(params.Title)
	params.CoverImageUrl = strings.TrimSpace(params.CoverImageUrl)
	if params.ID == uuid.Nil || params.ArtistID == uuid.Nil || params.Title == "" {
		return db.Album{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.UpdateAlbum(dbCtx, params)
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.DeleteAlbum(dbCtx, id)
}

func (s *AlbumService) ListAlbumSongs(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]db.Song, error) {
	if albumID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListAlbumSongs(dbCtx, albumID, limit, offset)
}
