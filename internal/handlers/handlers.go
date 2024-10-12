package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fallrising/goku-api/internal/database"
	"github.com/fallrising/goku-api/internal/models"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	db *database.Database
}

func NewBookmarkHandler(db *database.Database) *BookmarkHandler {
	return &BookmarkHandler{db: db}
}

func (h *BookmarkHandler) HandleUpload(c *gin.Context) {
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

		existingBookmark, err := h.db.GetBookmarkByURL(info.URL)
		if err != nil {
			errors = append(errors, "Error checking existing bookmark: "+err.Error())
			continue
		}

		if existingBookmark != nil {
			errors = append(errors, "Bookmark already exists: "+info.URL)
			continue
		}

		if err := h.db.SaveBookmark(info); err != nil {
			errors = append(errors, "Failed to save bookmark: "+err.Error())
			continue
		}

		processedURLs = append(processedURLs, info.URL)
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

func (h *BookmarkHandler) HandleGetAll(c *gin.Context) {
	bookmarks, err := h.db.GetAllBookmarks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookmarks"})
		return
	}

	c.JSON(http.StatusOK, bookmarks)
}

func (h *BookmarkHandler) HandleGetByURL(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	bookmark, err := h.db.GetBookmarkByURL(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookmark"})
		return
	}

	if bookmark == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, bookmark)
}

func (h *BookmarkHandler) HandleUpdate(c *gin.Context) {
	var info models.URLInfo
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := validateURLInfo(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateBookmark(info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark updated successfully"})
}

func (h *BookmarkHandler) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	if err := h.db.DeleteBookmark(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
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
