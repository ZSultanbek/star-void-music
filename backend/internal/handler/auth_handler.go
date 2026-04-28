package handler

import (
	"errors"
	"net/http"

	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func RegisterAuthRoutes(api *gin.RouterGroup, h *AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	result, err := h.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid credentials")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusCreated, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := bindJSONStrict(c, &req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	result, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			errorResponse(c, http.StatusBadRequest, "invalid credentials")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			errorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	dataResponse(c, http.StatusOK, result)
}

