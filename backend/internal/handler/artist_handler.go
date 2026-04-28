package handler

import (
	"errors"
	"net/http"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	service *service.ArtistService
}

func NewArtistHandler(s *service.ArtistService) *ArtistHandler {
	return &ArtistHandler{service: s}
}

func RegisterArtistRoutes(api *gin.RouterGroup, h *ArtistHandler, authMW gin.HandlerFunc, requireAdminMW gin.HandlerFunc) {
	artists := api.Group("/artists", authMW)
	{
		artists.GET("", h.ListArtists)
		artists.GET("/:id/albums", h.ListArtistAlbums)
		artists.GET("/:id", h.GetArtistByID)
	
		//admin routes
		artists.POST("", requireAdminMW,  h.CreateArtist)
		artists.PATCH("/:id", requireAdminMW, h.UpdateArtist)
		artists.DELETE("/:id", requireAdminMW, h.DeleteArtist)
	}
}

type createArtistRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type updateArtistRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (h *ArtistHandler) CreateArtist(c *gin.Context) {
	var req createArtistRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	artist, err := h.service.CreateArtist(ctx, service.CreateArtistInput{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid artist data")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusCreated, artist)
}

func (h *ArtistHandler) GetArtistByID(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid artist id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	artist, err := h.service.GetArtistByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid artist id")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "artist not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, artist)
}

func (h *ArtistHandler) UpdateArtist(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid artist id")
		return
	}

	var req updateArtistRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	artist, err := h.service.UpdateArtist(ctx, service.UpdateArtistInput{
		ID:   id,
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid artist data")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "artist not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, artist)
}

func (h *ArtistHandler) DeleteArtist(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid artist id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	if err := h.service.DeleteArtist(ctx, id); err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid artist id")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "artist not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, gin.H{"removed": true})
}

func (h *ArtistHandler) ListArtists(c *gin.Context) {
	limit, offset, err := parsePagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	artists, err := h.service.ListArtists(ctx, limit, offset)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, artists)
}

func (h *ArtistHandler) ListArtistAlbums(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid artist id")
		return
	}

	limit, offset, err := parsePagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	albums, err := h.service.ListArtistAlbums(ctx, id, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid artist id")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, albums)
}
