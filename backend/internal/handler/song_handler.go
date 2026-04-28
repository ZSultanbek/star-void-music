package handler

import (
	"errors"
	"io"
	"net/http"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SongHandler struct {
	service *service.SongService
}

func NewSongHandler(s *service.SongService) *SongHandler {
	return &SongHandler{service: s}
}

func RegisterSongRoutes(api *gin.RouterGroup, h *SongHandler, authMW gin.HandlerFunc, requireAdminMW gin.HandlerFunc) {
	songs := api.Group("/songs", authMW)
	{
		songs.GET("/:id", h.GetSongByID)
		songs.GET("/:id/stream", h.StreamSong)

		//admin routes
		songs.POST("", requireAdminMW, h.CreateSong)
		songs.PATCH("/:id", requireAdminMW, h.UpdateSong)
		songs.DELETE("/:id", requireAdminMW, h.DeleteSong)
	}
}

type createSongRequest struct {
	Title    string `json:"title"`
	AlbumID  string `json:"album_id"`
	Filepath string `json:"filepath"`
	Duration int32  `json:"duration"`
}

type updateSongRequest struct {
	Title    string `json:"title"`
	AlbumID  string `json:"album_id"`
	Filepath string `json:"filepath"`
	Duration int32  `json:"duration"`
}

func (h *SongHandler) CreateSong(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createSongRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	albumID, err := uuid.Parse(req.AlbumID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid album id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	song, err := h.service.CreateSong(ctx, service.CreateSongInput{
		Title:      req.Title,
		AlbumID:    albumID,
		Filepath:   req.Filepath,
		Duration:   req.Duration,
		UploadedBy: userID,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid song data")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusCreated, song)
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid song id")
		return
	}

	var req updateSongRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	albumID, err := uuid.Parse(req.AlbumID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid album id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	song, err := h.service.UpdateSong(ctx, service.UpdateSongInput{
		ID:       id,
		Title:    req.Title,
		AlbumID:  albumID,
		Filepath: req.Filepath,
		Duration: req.Duration,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid song data")
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

func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid song id")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	if err := h.service.DeleteSong(ctx, id); err != nil {
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

	dataResponse(c, http.StatusOK, gin.H{"removed": true})
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

	stream, err := h.service.GetStream(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			errorResponse(c, http.StatusNotFound, "song not found")
		case errors.Is(err, service.ErrFileNotFound):
			errorResponse(c, http.StatusNotFound, "song file not found")
		default:
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	defer func() {
		if closer, ok := stream.Reader.(io.Closer); ok {
			closer.Close()
		}
	}()

	c.Header("Content-Type", stream.ContentType)
	http.ServeContent(c.Writer, c.Request, stream.Name, stream.ModTime, stream.Reader)
}
