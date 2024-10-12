package handlers

import (
	"net/http"
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

// HandleUpload handles the upload of new bookmarks
func (h *BookmarkHandler) HandleUpload(c *gin.Context) {
	var urlInfos []models.URLInfo

	if err := c.BindJSON(&urlInfos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Rest of the HandleUpload function remains the same
	// ...
}

// HandleGetAll retrieves all bookmarks
func (h *BookmarkHandler) HandleGetAll(c *gin.Context) {
	bookmarks, err := h.db.GetAllBookmarks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookmarks"})
		return
	}

	c.JSON(http.StatusOK, bookmarks)
}

// HandleGetByURL retrieves a bookmark by URL
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

// HandleUpdate updates an existing bookmark
func (h *BookmarkHandler) HandleUpdate(c *gin.Context) {
	var info models.URLInfo
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.db.UpdateBookmark(info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark updated successfully"})
}

// HandleDelete deletes a bookmark by ID
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
