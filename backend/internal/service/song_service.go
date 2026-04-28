package service

import (
	"context"
	"errors"
	"strings"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SongRepository interface {
	CreateSong(ctx context.Context, params db.CreateSongParams) (db.Song, error)
	GetSongByID(ctx context.Context, id uuid.UUID) (db.Song, error)
	ListSongs(ctx context.Context, limit, offset int32) ([]db.Song, error)
	ListSongsByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]db.Song, error)
	UpdateSong(ctx context.Context, params db.UpdateSongParams) (db.Song, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
}

type SongService struct {
	repo SongRepository
}

func NewSongService(repo SongRepository) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) CreateSong(ctx context.Context, params db.CreateSongParams) (db.Song, error) {
	params.Title = strings.TrimSpace(params.Title)
	params.Filepath = strings.TrimSpace(params.Filepath)
	if params.Title == "" || params.Filepath == "" || params.Duration < 0 || params.AlbumID == uuid.Nil || params.UploadedBy == uuid.Nil {
		return db.Song{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.CreateSong(dbCtx, params)
}

func (s *SongService) GetSongByID(ctx context.Context, id uuid.UUID) (db.Song, error) {
	if id == uuid.Nil {
		return db.Song{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	song, err := s.repo.GetSongByID(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return db.Song{}, ErrNotFound
	}
	return song, err
}

func (s *SongService) ListSongs(ctx context.Context, limit, offset int32) ([]db.Song, error) {
	limit, offset = normalizePagination(limit, offset)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListSongs(dbCtx, limit, offset)
}

func (s *SongService) ListSongsByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]db.Song, error) {
	if albumID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListSongsByAlbumID(dbCtx, albumID, limit, offset)
}

func (s *SongService) UpdateSong(ctx context.Context, params db.UpdateSongParams) (db.Song, error) {
	params.Title = strings.TrimSpace(params.Title)
	params.Filepath = strings.TrimSpace(params.Filepath)
	if params.ID == uuid.Nil || params.AlbumID == uuid.Nil || params.Title == "" || params.Filepath == "" || params.Duration < 0 {
		return db.Song{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.UpdateSong(dbCtx, params)
}

func (s *SongService) DeleteSong(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.DeleteSong(dbCtx, id)
}
