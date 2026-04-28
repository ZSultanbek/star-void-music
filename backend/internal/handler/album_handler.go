package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlbumHandler struct {
	service *service.AlbumService
}

func NewAlbumHandler(s *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: s}
}

func RegisterAlbumRoutes(api *gin.RouterGroup, h *AlbumHandler, authMW gin.HandlerFunc, requireAdminMW gin.HandlerFunc) {
	albums := api.Group("/albums", authMW)
	{
		albums.GET("/:id", h.GetAlbumByID)
		albums.GET("/:id/songs", h.GetAlbumSongs)

		//admin routes
		albums.POST("", requireAdminMW, h.CreateAlbum)
		albums.PATCH("/:id", requireAdminMW, h.UpdateAlbum)
		albums.DELETE("/:id", requireAdminMW, h.DeleteAlbum)
	}
}

type createAlbumRequest struct {
	Title         string  `json:"title"`
	ArtistID      string  `json:"artist_id"`
	CoverImageURL string  `json:"cover_image_url"`
	ReleaseDate   *string `json:"release_date"`
}

type updateAlbumRequest struct {
	Title         string  `json:"title"`
	ArtistID      string  `json:"artist_id"`
	CoverImageURL string  `json:"cover_image_url"`
	ReleaseDate   *string `json:"release_date"`
}

func parseReleaseDate(in *string) (*time.Time, error) {
	if in == nil || strings.TrimSpace(*in) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*in))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	var req createAlbumRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	artistID, err := uuid.Parse(req.ArtistID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid artist id")
		return
	}

	releaseDate, err := parseReleaseDate(req.ReleaseDate)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid release_date")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	album, err := h.service.CreateAlbum(ctx, service.CreateAlbumInput{
		Title:         req.Title,
		ArtistID:      artistID,
		CoverImageURL: req.CoverImageURL,
		ReleaseDate:   releaseDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid album data")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusCreated, album)
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid album id")
		return
	}

	var req updateAlbumRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	artistID, err := uuid.Parse(req.ArtistID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid artist id")
		return
	}

	releaseDate, err := parseReleaseDate(req.ReleaseDate)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid release_date")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	album, err := h.service.UpdateAlbum(ctx, service.UpdateAlbumInput{
		ID:            id,
		Title:         req.Title,
		ArtistID:      artistID,
		CoverImageURL: req.CoverImageURL,
		ReleaseDate:   releaseDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid album data")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "album not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, album)
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid album id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	if err := h.service.DeleteAlbum(ctx, id); err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid album id")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "album not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, gin.H{"removed": true})
}

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid album id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	album, err := h.service.GetAlbumByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid album id")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "album not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, album)
}

func (h *AlbumHandler) GetAlbumSongs(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid album id")
		return
	}

	limit, offset, err := parsePagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	songs, err := h.service.ListAlbumSongs(ctx, id, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid album id")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, songs)
}
