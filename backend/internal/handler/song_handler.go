package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	service *service.SongService
}

func NewSongHandler(s *service.SongService) *SongHandler {
	return &SongHandler{service: s}
}

func RegisterSongRoutes(api *gin.RouterGroup, h *SongHandler, authMW gin.HandlerFunc) {
	songs := api.Group("/songs", authMW)
	{
		songs.GET("/:id", h.GetSongByID)
		songs.GET("/:id/stream", h.StreamSong)
	}
}

func (h *SongHandler) GetSongByID(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid song id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	song, err := h.service.GetSongByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid song id")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "song not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, song)
}

func (h *SongHandler) StreamSong(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid song id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	song, err := h.service.GetSongByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "song not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	file, err := os.Open(song.Filepath)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "song file not found")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		errorResponse(c, http.StatusNotFound, "song file not found")
		return
	}

	c.Header("Content-Type", "audio/mpeg")
	http.ServeContent(c.Writer, c.Request, filepath.Base(song.Filepath), stat.ModTime(), file)
}
