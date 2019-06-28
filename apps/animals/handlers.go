package animals

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qwezarty/zoo-demo/apps"
	"github.com/qwezarty/zoo-demo/models"
)

type AnimalAPIs struct {
	apps.RestAPIs
}

type AminalViewModel struct {
	Base    models.Base
	Animals []models.Animal
}

func (a *AnimalAPIs) Hello(c *gin.Context) {
	c.String(http.StatusOK, "hello from an animal!")
}

func (a *AnimalAPIs) Create(c *gin.Context) {
	animal := &models.Animal{}
	zoo := &models.Zoo{}
	if err := c.Bind(animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if zoo existed
	db.Where("id = ?", animal.ZooID).First(zoo)

	// create animal if zoo is already existed
	if zoo.ID != "" && len(zoo.ID) != 36 { // length of uuid is 36
		if err := db.Create(animal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, animal)
		return
	}

	// else begin a transaction to create both zoo and animal
	zoo.Name = "unnamed"
	tx := db.Begin()
	if err := db.Create(zoo).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	animal.ZooID = zoo.ID
	if err := db.Create(animal).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, animal)
	return
}
