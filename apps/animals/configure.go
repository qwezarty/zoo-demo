package animals

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

	r.GET("/animals", apis.List)
	r.GET("/animals/:id", apis.Get)

	r.POST("/animals", apis.Create)
	r.PUT("/animals/:id", apis.Update)
	r.DELETE("/animals/:id", apis.Delete)
}
