package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"url_shortener/models"
	"url_shortener/repositories"
	"url_shortener/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	minHashLength = 7
)

type URLController struct {
	repo    repositories.URLRepository
	baseURL string
	salt    string
}

// NewURLController creates a new URLController instance with injected URLRepository.
func NewURLController(repo repositories.URLRepository, baseURL string, salt string) *URLController {
	return &URLController{
		repo:    repo,
		baseURL: baseURL,
		salt:    salt,
	}
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

// ShortenURL handles shortening requests with the following logic:
// 1. Check if the long URL already exists in the DB; if yes, return its short URL.
// 2. If not, attempt to insert a new record with a hash of the long URL.
// 3. On collision (hash exists for a different URL), append a fixed salt (from environment)
//    to the long URL and retry, possibly increasing hash length after several attempts.
func (uc *URLController) ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	longURL := strings.TrimSpace(req.URL)
	if longURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL cannot be empty"})
		return
	}

	// 1. Check if long URL already exists
	existing, err := uc.repo.FindByLongURL(longURL)
	if err == nil && existing != nil {
		// URL already shortened; return existing short URL
		c.JSON(http.StatusOK, gin.H{
			"url":       existing.LongURL,
			"short_url": uc.baseURL + "/" + existing.ShortURL,
		})
		return
	}
	// If err != nil and not record not found, treat as error (but we'll assume not found for now)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check URL existence: " + err.Error()})
		return
	}

	// 2. Attempt to insert new URL with collision handling using salt
	var shortURL string
	var success bool

	// Retry loop: start with length 7, try with progressively more salt characters on each retry.
	// We terminate when we've used all characters in the salt.
	maxAttempt := len(uc.salt)
	if maxAttempt == 0 {
		maxAttempt = 1 // Still try once even if salt is empty
	}
	for attempt := 0; attempt < maxAttempt; attempt++ {
		length := minHashLength + attempt
		input := longURL
		if attempt > 0 {
			// Append progressively more characters of salt for each retry after the first
			input = longURL + uc.salt[:attempt]
		}
		hash := utils.GenerateURLHash(input, length)

		newURL := models.URL{
			ShortURL: hash,
			LongURL:  longURL,
		}

		err := uc.repo.Create(&newURL)
		if err == nil {
			shortURL = hash
			success = true
			break
		}

		// Check if the insertion failed due to unique constraint conflict
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "UNIQUE") {
			// Collision: hash already exists for a different long URL (we already checked longURL doesn't exist)
			continue
		}

		// Any other DB error: abort and return internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL: " + err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to generate unique short URL after " + strconv.Itoa(maxAttempt) + " attempts due to collision"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":       longURL,
		"short_url": uc.baseURL + "/" + shortURL,
	})
}

// RedirectURL handles redirection from short URL to original long URL.
func (uc *URLController) RedirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short URL"})
		return
	}

	url, err := uc.repo.FindByShortURL(shortURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.LongURL)
}
