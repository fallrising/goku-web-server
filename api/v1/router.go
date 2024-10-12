package v1

import (
	"github.com/fallrising/goku-api/internal/database"
	"github.com/fallrising/goku-api/internal/handlers/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *database.Database) *gin.RouterGroup {
	router := gin.New()
	v1 := router.Group("/api/v1")

	bookmarkHandler := handlers.NewBookmarkHandler(db)

	v1.POST("/upload", bookmarkHandler.HandleUpload)
	v1.GET("/bookmarks", bookmarkHandler.HandleGetAll)
	v1.GET("/bookmark", bookmarkHandler.HandleGetByURL)
	v1.PUT("/bookmark", bookmarkHandler.HandleUpdate)
	v1.DELETE("/bookmark/:id", bookmarkHandler.HandleDelete)

	return v1
}
