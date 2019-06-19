package zoo

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Configure registers all handler
func Configure(r *gin.Engine, engine *gorm.DB) {
	db = engine

	r.GET("/zoo/:id", Get)
	r.GET("/zoo", List)

	r.POST("/zoo", Create)
	r.PUT("/zoo/:id", Update)
	r.DELETE("/zoo/:id", Delete)
}
