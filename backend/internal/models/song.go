package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	AlbumID    uuid.UUID `json:"album_id"`
	Filepath   string    `json:"filepath"`
	Duration   int32     `json:"duration"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateSongInput struct {
	Title      string    `json:"title"`
	AlbumID    uuid.UUID `json:"album_id"`
	Filepath   string    `json:"filepath"`
	Duration   int32     `json:"duration"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
}

type UpdateSongInput struct {
	Title    string    `json:"title"`
	AlbumID  uuid.UUID `json:"album_id"`
	Filepath string    `json:"filepath"`
	Duration int32     `json:"duration"`
}
