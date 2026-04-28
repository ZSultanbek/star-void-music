package service

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	db "star-void-music/backend/db/sqlc"
	"star-void-music/backend/internal/models"

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

type Stream struct {
	Reader      io.ReadSeeker
	Size        int64
	ContentType string
	Name        string
	ModTime     time.Time
}

var ErrFileNotFound = errors.New("file not found")

func NewSongService(repo SongRepository) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) CreateSong(ctx context.Context, in CreateSongInput) (models.Song, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Filepath = strings.TrimSpace(in.Filepath)
	if in.Title == "" || in.Filepath == "" || in.Duration < 0 || in.AlbumID == uuid.Nil || in.UploadedBy == uuid.Nil {
		return models.Song{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row, err := s.repo.CreateSong(dbCtx, db.CreateSongParams{
		Title:      in.Title,
		AlbumID:    in.AlbumID,
		Filepath:   in.Filepath,
		Duration:   in.Duration,
		UploadedBy: in.UploadedBy,
	})
	if err != nil {
		return models.Song{}, err
	}
	return mapSong(row), nil
}

func (s *SongService) GetSongByID(ctx context.Context, id uuid.UUID) (models.Song, error) {
	if id == uuid.Nil {
		return models.Song{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row, err := s.repo.GetSongByID(dbCtx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return models.Song{}, ErrNotFound
	}
	if err != nil {
		return models.Song{}, err
	}
	return mapSong(row), nil
}

func (s *SongService) ListSongs(ctx context.Context, limit, offset int32) ([]models.Song, error) {
	limit, offset = normalizePagination(limit, offset)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	rows, err := s.repo.ListSongs(dbCtx, limit, offset)
	if err != nil {
		return nil, err
	}
	return mapSongs(rows), nil
}

func (s *SongService) ListSongsByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int32) ([]models.Song, error) {
	if albumID == uuid.Nil {
		return nil, ErrValidation
	}
	limit, offset = normalizePagination(limit, offset)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	rows, err := s.repo.ListSongsByAlbumID(dbCtx, albumID, limit, offset)
	if err != nil {
		return nil, err
	}
	return mapSongs(rows), nil
}

func (s *SongService) UpdateSong(ctx context.Context, in UpdateSongInput) (models.Song, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Filepath = strings.TrimSpace(in.Filepath)
	if in.ID == uuid.Nil || in.AlbumID == uuid.Nil || in.Title == "" || in.Filepath == "" || in.Duration < 0 {
		return models.Song{}, ErrValidation
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	row, err := s.repo.UpdateSong(dbCtx, db.UpdateSongParams{
		ID:       in.ID,
		Title:    in.Title,
		AlbumID:  in.AlbumID,
		Filepath: in.Filepath,
		Duration: in.Duration,
	})
	if err != nil {
		return models.Song{}, err
	}
	return mapSong(row), nil
}

func (s *SongService) DeleteSong(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrValidation
	}
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return s.repo.DeleteSong(dbCtx, id)
}

func (s *SongService) GetStream(ctx context.Context, id uuid.UUID) (*Stream, error) {
	song, err := s.GetSongByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	file, err := os.Open(song.Filepath)
	if err != nil {
		return nil, ErrFileNotFound
	}

	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		return nil, ErrFileNotFound
	}

	return &Stream{
		Reader:      file,
		Size:        stat.Size(),
		ContentType: "audio/mpeg",
		Name:        filepath.Base(song.Filepath),
		ModTime:     stat.ModTime(),
	}, nil
}

func mapSong(row db.Song) models.Song {
	return models.Song{
		ID:         row.ID,
		Title:      row.Title,
		AlbumID:    row.AlbumID,
		Filepath:   row.Filepath,
		Duration:   row.Duration,
		UploadedBy: row.UploadedBy,
		CreatedAt:  row.CreatedAt,
	}
}

func mapSongs(rows []db.Song) []models.Song {
	out := make([]models.Song, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapSong(r))
	}
	return out
}

type CreateSongInput struct {
	Title      string
	AlbumID    uuid.UUID
	Filepath   string
	Duration   int32
	UploadedBy uuid.UUID
}

type UpdateSongInput struct {
	ID       uuid.UUID
	Title    string
	AlbumID  uuid.UUID
	Filepath string
	Duration int32
}
