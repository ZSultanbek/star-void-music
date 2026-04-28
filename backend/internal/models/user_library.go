package models

import (
	"time"

	"github.com/google/uuid"
)

type UserLibraryItem struct {
	UserID  uuid.UUID `json:"user_id"`
	SongID  uuid.UUID `json:"song_id"`
	AddedAt time.Time `json:"added_at"`
}

type AddSongToLibraryInput struct {
	UserID uuid.UUID `json:"user_id"`
	SongID uuid.UUID `json:"song_id"`
}

type RemoveSongFromLibraryInput struct {
	UserID uuid.UUID `json:"user_id"`
	SongID uuid.UUID `json:"song_id"`
}
