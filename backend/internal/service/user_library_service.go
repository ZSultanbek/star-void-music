package service

import (
	"context"

	db "star-void-music/backend/db/sqlc"

	"github.com/google/uuid"
)

type UserLibraryRepository interface {
	AddSongToUserLibrary(ctx context.Context, params db.AddSongToUserLibraryParams) (db.UserLibrary, error)
	ListUserLibrarySongs(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.ListUserLibrarySongsRow, error)
	RemoveSongFromUserLibrary(ctx context.Context, params db.RemoveSongFromUserLibraryParams) error
}

type UserLibraryService struct {
	repo UserLibraryRepository
}

func NewUserLibraryService(repo UserLibraryRepository) *UserLibraryService {
	return &UserLibraryService{repo: repo}
}

func (s *UserLibraryService) AddSongToLibrary(ctx context.Context, params db.AddSongToUserLibraryParams) (db.UserLibrary, error) {
	if params.UserID == uuid.Nil || params.SongID == uuid.Nil {
		return db.UserLibrary{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.AddSongToUserLibrary(dbCtx, params)
}

func (s *UserLibraryService) ListLibrarySongs(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.ListUserLibrarySongsRow, error) {
	if userID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.ListUserLibrarySongs(dbCtx, userID, limit, offset)
}

func (s *UserLibraryService) RemoveSongFromLibrary(ctx context.Context, params db.RemoveSongFromUserLibraryParams) error {
	if params.UserID == uuid.Nil || params.SongID == uuid.Nil {
		return ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.RemoveSongFromUserLibrary(dbCtx, params)
}
