package zoos

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func List(c *gin.Context) {
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

func Delete(c *gin.Context) {
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
