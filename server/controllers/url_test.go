package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockURLRepository is a mock implementation of repositories.URLRepository using testify mock.
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Create(url *models.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func TestURLController_ShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful shorten on first attempt (length 7)", func(t *testing.T) {
		mockRepo := new(MockURLRepository)

		// Set up mock expectation for length 7
		mockRepo.On("Create", mock.MatchedBy(func(u *models.URL) bool {
			return len(u.ShortURL) == 7 && u.LongURL == "https://example.com"
		})).Return(nil).Once()

		controller := NewURLController(mockRepo, "http://localhost:8080")

		r := gin.New()
		r.POST("/shorten", controller.ShortenURL)

		reqBody := []byte(`{"url": "https://example.com"}`)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response body: %v", err)
		}

		if resp["url"] != "https://example.com" {
			t.Errorf("expected url 'https://example.com', got '%s'", resp["url"])
		}

		shortURL := resp["short_url"]
		expectedPrefix := "http://localhost:8080/"
		if len(shortURL) != len(expectedPrefix)+7 {
			t.Errorf("expected short_url hash length to be 7, got short_url: %s", shortURL)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("shorten succeeds on retry (first fails with duplicate, second passes, length 8)", func(t *testing.T) {
		mockRepo := new(MockURLRepository)

		// First call (length 7) returns duplicate key error
		mockRepo.On("Create", mock.MatchedBy(func(u *models.URL) bool {
			return len(u.ShortURL) == 7 && u.LongURL == "https://example.com"
		})).Return(gorm.ErrDuplicatedKey).Once()

		// Second call (length 8) returns success (nil)
		mockRepo.On("Create", mock.MatchedBy(func(u *models.URL) bool {
			return len(u.ShortURL) == 8 && u.LongURL == "https://example.com"
		})).Return(nil).Once()

		controller := NewURLController(mockRepo, "http://localhost:8080")

		r := gin.New()
		r.POST("/shorten", controller.ShortenURL)

		reqBody := []byte(`{"url": "https://example.com"}`)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		shortURL := resp["short_url"]

		expectedPrefix := "http://localhost:8080/"
		if len(shortURL) != len(expectedPrefix)+8 {
			t.Errorf("expected short_url hash length to be 8, got short_url: %s", shortURL)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("conflict error when all 4 attempts fail due to collision", func(t *testing.T) {
		mockRepo := new(MockURLRepository)

		// All 4 attempts (lengths 7, 8, 9, 10) return duplicate key error
		for i := 0; i < 4; i++ {
			length := 7 + i
			mockRepo.On("Create", mock.MatchedBy(func(u *models.URL) bool {
				return len(u.ShortURL) == length && u.LongURL == "https://example.com"
			})).Return(gorm.ErrDuplicatedKey).Once()
		}

		controller := NewURLController(mockRepo, "http://localhost:8080")

		r := gin.New()
		r.POST("/shorten", controller.ShortenURL)

		reqBody := []byte(`{"url": "https://example.com"}`)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status 409 Conflict, got %d", w.Code)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("internal server error immediately on non-collision DB error", func(t *testing.T) {
		mockRepo := new(MockURLRepository)

		// First call returns a generic database error and aborts immediately
		mockRepo.On("Create", mock.MatchedBy(func(u *models.URL) bool {
			return len(u.ShortURL) == 7 && u.LongURL == "https://example.com"
		})).Return(errors.New("connection reset by peer")).Once()

		controller := NewURLController(mockRepo, "http://localhost:8080")

		r := gin.New()
		r.POST("/shorten", controller.ShortenURL)

		reqBody := []byte(`{"url": "https://example.com"}`)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500 Internal Server Error, got %d", w.Code)
		}

		mockRepo.AssertExpectations(t)
	})
}
