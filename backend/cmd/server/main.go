package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"star-void-music/backend/internal/handler"
	"star-void-music/backend/internal/middleware"
	"star-void-music/backend/internal/repository"
	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := repository.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService, err := service.NewAuthServiceFromEnv(userService, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	authHandler := handler.NewAuthHandler(authService)

	authMW, err := middleware.NewJWTAuthMiddlewareFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	requireAdminMW := middleware.RequireRole("admin")

	artistRepo := repository.NewArtistRepository(pool)
	artistService := service.NewArtistService(artistRepo)
	artistHandler := handler.NewArtistHandler(artistService)

	albumRepo := repository.NewAlbumRepository(pool)
	albumService := service.NewAlbumService(albumRepo)
	albumHandler := handler.NewAlbumHandler(albumService)

	songRepo := repository.NewSongRepository(pool)
	songService := service.NewSongService(songRepo)
	songHandler := handler.NewSongHandler(songService)

	userLibraryRepo := repository.NewUserLibraryRepository(pool)
	userLibraryService := service.NewUserLibraryService(userLibraryRepo)
	userLibraryHandler := handler.NewUserLibraryHandler(userLibraryService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "star-void-music-backend",
		})
	})

	api := r.Group("/api")
	handler.RegisterAuthRoutes(api, authHandler)
	handler.RegisterUserRoutes(api, userHandler, authMW)
	handler.RegisterArtistRoutes(api, artistHandler, authMW, requireAdminMW)
	handler.RegisterAlbumRoutes(api, albumHandler, authMW, requireAdminMW)
	handler.RegisterSongRoutes(api, songHandler, authMW, requireAdminMW)
	handler.RegisterUserLibraryRoutes(api, userLibraryHandler, authMW)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}