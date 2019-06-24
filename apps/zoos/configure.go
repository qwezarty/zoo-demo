package zoos

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Configure registers all handler
func Configure(r *gin.Engine, engine *gorm.DB) {
	db = engine

	r.GET("/zoos/:id", Get)
	r.GET("/zoos", List)

	r.POST("/zoos", Create)
	r.PUT("/zoos/:id", Update)
	r.DELETE("/zoos/:id", Delete)
}
