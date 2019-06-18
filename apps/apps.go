package apps

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func Configure(engine *gorm.DB) {
	db = engine
}

type RestAPIs struct {
	// Model interface{}
	Bean  interface{}
	Beans interface{}
}

func (a *RestAPIs) Gets(c *gin.Context) {
	db.Find(a.Beans)
	c.JSON(http.StatusOK, a.Beans)
}

func (a *RestAPIs) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// using reflect to create
	// bean := reflect.New(a.Model)
	// db.Where("id = ?", id).First(bean.Interface())
	//c.JSON(http.StatusOK, bean.Elem().Interface())

	// using beans to create
	db.Where("id = ?", id).First(a.Bean)
	c.JSON(http.StatusOK, a.Bean)
}

func (a *RestAPIs) Create(c *gin.Context) {
	if err := c.Bind(a.Bean); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// using reflect to set id field
	u, _ := uuid.NewUUID()
	rv := reflect.ValueOf(a.Bean).Elem()
	rv.FieldByName("ID").SetString(u.String())

	db.Create(a.Bean)
	c.JSON(http.StatusOK, a.Bean)
}

func (a *RestAPIs) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	if err := c.Bind(a.Bean); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(a.Bean).Updates(a.Bean)
	c.JSON(http.StatusOK, a.Bean)
}

func (a *RestAPIs) Remove(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	rv := reflect.ValueOf(a.Bean).Elem()
	rv.FieldByName("ID").SetString(id)

	// soft delete
	// db.Delete(&bean)

	// delete permanently
	db.Unscoped().Delete(a.Bean)
}
