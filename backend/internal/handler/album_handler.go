package handler

import (
	"errors"
	"net/http"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	service *service.AlbumService
}

func NewAlbumHandler(s *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: s}
}

func RegisterAlbumRoutes(api *gin.RouterGroup, h *AlbumHandler, authMW gin.HandlerFunc) {
	albums := api.Group("/albums", authMW)
	{
		albums.GET("/:id", h.GetAlbumByID)
		albums.GET("/:id/songs", h.GetAlbumSongs)
	}
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
