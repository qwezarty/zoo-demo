package animal

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qwezarty/zoo-demo/models"
)

var db *gorm.DB

// Configure registers all handler
func Configure(r *gin.Engine, engine *gorm.DB) {
	db = engine
	apis := AnimalAPIs{}
	apis.Bean = &models.Animal{}
	apis.Beans = &[]models.Animal{}

	r.GET("/animal", apis.List)
	r.GET("/animal/:id", apis.Get)

	r.POST("/animal", apis.Create)
	r.PUT("/animal/:id", apis.Update)
	r.DELETE("/animal/:id", apis.Delete)
}
