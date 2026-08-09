package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"url_shortener/controllers"
	"url_shortener/db"
	"url_shortener/repositories"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	database, err := db.InitDB()
	if err != nil {
		os.Exit(1)
	}

	repo := repositories.NewSQLURLRepository(database)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	salt := os.Getenv("SHORTEN_SALT")
	if salt == "" {
		salt = "fixed_salt_change_in_production" // default for development
	}
	controller := controllers.NewURLController(repo, baseURL, salt)
	healthController := controllers.NewHealthController()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.Printf("Recovery from panic: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}))
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowOrigins []string
	if corsOrigins == "" {
		allowOrigins = []string{"http://localhost"}
	} else {
		parts := strings.Split(corsOrigins, ",")
		allowOrigins = make([]string, 0, len(parts))
		for _, p := range parts {
			allowOrigins = append(allowOrigins, strings.TrimSpace(p))
		}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})
	r.GET("/health", healthController.HealthCheck)
	r.POST("/api/v1/shorten", controller.ShortenURL)
	r.GET("/:shortURL", controller.RedirectURL)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run the server in a separate goroutine
	go func() {
		log.Println("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for OS interrupt signal
	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
