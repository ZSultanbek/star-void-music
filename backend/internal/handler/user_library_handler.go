package handler

import (
	"errors"
	"net/http"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserLibraryHandler struct {
	service *service.UserLibraryService
}

type addSongToLibraryRequest struct {
	SongID string `json:"song_id"`
}

func NewUserLibraryHandler(s *service.UserLibraryService) *UserLibraryHandler {
	return &UserLibraryHandler{service: s}
}

func RegisterUserLibraryRoutes(api *gin.RouterGroup, h *UserLibraryHandler, authMW gin.HandlerFunc) {
	me := api.Group("/me", authMW)
	{
		me.POST("/library", h.AddSongToLibrary)
		me.GET("/library", h.ListLibrarySongs)
		me.DELETE("/library/:song_id", h.RemoveSongFromLibrary)
	}
}

func (h *UserLibraryHandler) AddSongToLibrary(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req addSongToLibraryRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	songID, err := parseUUIDString(req.SongID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid song id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	item, err := h.service.AddSongToLibrary(ctx, userID, songID)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid library data")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusCreated, item)
}

func (h *UserLibraryHandler) ListLibrarySongs(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit, offset, err := parsePagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	items, err := h.service.ListLibrarySongs(ctx, userID, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, items)
}

func (h *UserLibraryHandler) RemoveSongFromLibrary(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	songID, err := parseUUIDParam(c, "song_id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid song id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	err = h.service.RemoveSongFromLibrary(ctx, userID, songID)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, gin.H{"removed": true})
}

func parseUUIDString(v string) (uuid.UUID, error) {
	return uuid.Parse(v)
}
