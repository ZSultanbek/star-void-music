package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	ArtistID      uuid.UUID  `json:"artist_id"`
	CoverImageURL string     `json:"cover_image_url"`
	ReleaseDate   *time.Time `json:"release_date,omitempty"` // date-only in DB
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateAlbumInput struct {
	Title         string     `json:"title"`
	ArtistID      uuid.UUID  `json:"artist_id"`
	CoverImageURL string     `json:"cover_image_url"`
	ReleaseDate   *time.Time `json:"release_date,omitempty"`
}

type UpdateAlbumInput struct {
	Title         string     `json:"title"`
	ArtistID      uuid.UUID  `json:"artist_id"`
	CoverImageURL string     `json:"cover_image_url"`
	ReleaseDate   *time.Time `json:"release_date,omitempty"`
}
