package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/fallrising/goku-api/internal/database"
	"github.com/fallrising/goku-api/internal/models"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	db *database.Database
}

func NewUploadHandler(db *database.Database) *UploadHandler {
	return &UploadHandler{db: db}
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	var urlInfos []models.URLInfo

	if err := c.BindJSON(&urlInfos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if len(urlInfos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No URL information provided"})
		return
	}

	var processedURLs []string
	var errors []string

	for _, info := range urlInfos {
		if err := validateURLInfo(info); err != nil {
			errors = append(errors, err.Error())
			continue
		}

		// Check if bookmark already exists
		existingBookmark, err := h.db.GetBookmarkByURL(info.URL)
		if err != nil {
			errors = append(errors, "Error checking existing bookmark: "+err.Error())
			continue
		}

		if existingBookmark != nil {
			errors = append(errors, "Bookmark already exists: "+info.URL)
			continue
		}

		// Save new bookmark
		if err := h.db.SaveBookmark(info); err != nil {
			errors = append(errors, "Failed to save bookmark: "+err.Error())
			continue
		}

		processedURLs = append(processedURLs, info.URL)
		// Use gin.DefaultWriter for logging
		gin.DefaultWriter.Write([]byte("Processed URL: " + info.URL + ", Title: " + info.Title + "\n"))
	}

	response := gin.H{
		"processed_urls": processedURLs,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	statusCode := http.StatusOK
	if len(processedURLs) == 0 && len(errors) > 0 {
		statusCode = http.StatusBadRequest
	} else if len(errors) > 0 {
		statusCode = http.StatusPartialContent
	}

	c.JSON(statusCode, response)
}

func validateURLInfo(info models.URLInfo) error {
	if info.URL == "" {
		return errors.New("URL is required")
	}

	if _, err := url.ParseRequestURI(info.URL); err != nil {
		return errors.New("Invalid URL format")
	}

	if info.Title == "" {
		return errors.New("Title is required")
	}

	return nil
}
