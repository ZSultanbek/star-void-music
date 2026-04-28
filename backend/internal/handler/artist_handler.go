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

func RegisterArtistRoutes(api *gin.RouterGroup, h *ArtistHandler, authMW gin.HandlerFunc) {
	artists := api.Group("/artists", authMW)
	{
		artists.GET("", h.ListArtists)
		artists.GET("/:id/albums", h.ListArtistAlbums)
	}
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
