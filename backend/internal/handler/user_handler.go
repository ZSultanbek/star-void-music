package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestTimeout = 3 * time.Second

type UserHandler struct {
	service *service.UserService
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func RegisterUserRoutes(api *gin.RouterGroup, h *UserHandler, authMW gin.HandlerFunc) {
	users := api.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", authMW, h.GetUserByID)
	}
	me := api.Group("/me", authMW)
	{
		me.GET("", h.GetMe)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	user, err := h.service.CreateUser(ctx, service.CreateUserInput{
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: req.Password,
		Role:         strings.TrimSpace(req.Role),
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid user data")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	targetID, err := parseUUIDParam(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	requesterID, ok := authUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	if requesterID != targetID && authUserRole(c) != "admin" {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	user, err := h.service.GetUserByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid user id")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, user)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, user)
}

func withTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), requestTimeout)
}

func bindJSONStrict(c *gin.Context, dst any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("unexpected trailing content")
	}
	return nil
}

func parseUUIDParam(c *gin.Context, key string) (uuid.UUID, error) {
	return uuid.Parse(strings.TrimSpace(c.Param(key)))
}

func parsePagination(c *gin.Context) (int32, int32, error) {
	limit := int32(20)
	offset := int32(0)

	if v := strings.TrimSpace(c.Query("limit")); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, 0, errors.New("invalid limit")
		}
		limit = int32(n)
	}
	if v := strings.TrimSpace(c.Query("offset")); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, 0, errors.New("invalid offset")
		}
		offset = int32(n)
	}
	return limit, offset, nil
}

func authUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return uuid.UUID{}, false
	}
	switch t := v.(type) {
	case uuid.UUID:
		if t == uuid.Nil {
			return uuid.UUID{}, false
		}
		return t, true
	case string:
		id, err := uuid.Parse(t)
		if err != nil || id == uuid.Nil {
			return uuid.UUID{}, false
		}
		return id, true
	default:
		return uuid.UUID{}, false
	}
}

func authUserRole(c *gin.Context) string {
	v, ok := c.Get("role")
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return strings.ToLower(strings.TrimSpace(s))
}

func dataResponse(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"data":  data,
		"error": nil,
	})
}

func errorResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"data":  nil,
		"error": msg,
	})
}