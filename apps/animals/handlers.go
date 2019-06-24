package animals

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qwezarty/zoo-demo/apps"
	"github.com/qwezarty/zoo-demo/models"
)

type AnimalAPIs struct {
	apps.RestAPIs
}

func (a *AnimalAPIs) Create(c *gin.Context) {
	animal := &models.Animal{}
	zoo := &models.Zoo{}
	if err := c.Bind(animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if zoo existed
	u, _ := uuid.NewUUID()
	animal.ID = u.String()
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
	u, _ = uuid.NewUUID()
	zoo.ID = u.String()
	zoo.Name = "unnamed"
	animal.ZooID = zoo.ID
	tx := db.Begin()
	if err := db.Create(zoo).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(animal).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, animal)
	return
}
