package models

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateArtistInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateArtistInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
