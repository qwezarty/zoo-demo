package zoo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qwezarty/zoo-demo/models"
)

func Get(c *gin.Context) {
	var bean models.Zoo
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	db.Where("id = ?", id).First(&bean)
	c.JSON(http.StatusOK, bean)
}

func Gets(c *gin.Context) {
	var beans = make([]models.Zoo, 0)
	db.Find(&beans)
	c.JSON(http.StatusOK, beans)
}

func Create(c *gin.Context) {
	var bean models.Zoo
	if err := c.Bind(&bean); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, _ := uuid.NewUUID()
	bean.ID = u.String()

	db.Create(&bean)
	c.JSON(http.StatusOK, bean)
}

func Update(c *gin.Context) {
	var bean models.Zoo
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	if err := c.Bind(&bean); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&bean).Updates(bean)
	c.JSON(http.StatusOK, bean)
}

func Remove(c *gin.Context) {
	var bean models.Zoo
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	bean.ID = id

	// soft delete
	// db.Delete(&bean)

	// delete permanently
	db.Unscoped().Delete(&bean)
}
