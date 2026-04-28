package service

import (
	"context"

	db "star-void-music/backend/db/sqlc"
	"star-void-music/backend/internal/models"

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

func (s *UserLibraryService) AddSongToLibrary(ctx context.Context, userID, songID uuid.UUID) (models.UserLibraryItem, error) {
	if userID == uuid.Nil || songID == uuid.Nil {
		return models.UserLibraryItem{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row, err := s.repo.AddSongToUserLibrary(dbCtx, db.AddSongToUserLibraryParams{
		UserID: userID,
		SongID: songID,
	})
	if err != nil {
		return models.UserLibraryItem{}, err
	}
	return models.UserLibraryItem{
		UserID: row.UserID,
		SongID: row.SongID,
		AddedAt: row.AddedAt,
	}, nil
}

func (s *UserLibraryService) ListLibrarySongs(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]models.Song, error) {
	if userID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := s.repo.ListUserLibrarySongs(dbCtx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	out := make([]models.Song, 0, len(rows))
	for _, r := range rows {
		out = append(out, models.Song{
			ID: r.SongID, Title: r.Title, AlbumID: r.AlbumID, Filepath: r.Filepath,
			Duration: r.Duration, UploadedBy: r.UploadedBy, CreatedAt: r.CreatedAt,
		})
	}
	return out, nil
}

func (s *UserLibraryService) RemoveSongFromLibrary(ctx context.Context, userID, songID uuid.UUID) error {
	if userID == uuid.Nil || songID == uuid.Nil {
		return ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	return s.repo.RemoveSongFromUserLibrary(dbCtx, db.RemoveSongFromUserLibraryParams{
		UserID: userID,
		SongID: songID,
	})
}
