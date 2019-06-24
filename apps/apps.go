package apps

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qwezarty/zoo-demo/models"
)

var db *gorm.DB

func Configure(engine *gorm.DB) {
	db = engine
}

type RestAPIs struct {
	// model and models
	Bean  interface{}
	Beans interface{}
}

func (a *RestAPIs) List(c *gin.Context) {
	q := db.New()
	if got := c.Query("start_time"); got != "" {
		t, err := time.Parse(time.RFC3339, got)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("error when parsing time, RFC3339 and url encode is needed: %v", err),
			})
			return
		}
		q = q.Where("created_at > ?", t)
	}
	if got := c.Query("end_time"); got != "" {
		t, err := time.Parse(time.RFC3339, got)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("error when parsing time, RFC3339 and url encode is needed: %v", err),
			})
			return
		}
		// qcount = qcount.Where("created_at < ?", t)
		q = q.Where("created_at < ?", t)
	}
	// get pagination infos from query string
	count, ptoken, psize := 0, 0, 10
	if got := c.Query("page_token"); got != "" {
		num, err := strconv.Atoi(got)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_token should be a number"})
			return
		}
		ptoken = num
	}
	if got := c.Query("page_size"); got != "" {
		num, err := strconv.Atoi(got)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size should be a number"})
			return
		}
		psize = num
	}

	if err := q.Find(a.Beans).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := q.Offset(ptoken * psize).Limit(psize).Find(a.Beans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ret := models.Pagination{
		PageToken: ptoken,
		PageSize:  psize,
		TotalSize: count,
		Body:      a.Beans,
	}

	c.JSON(http.StatusOK, ret)
}

func (a *RestAPIs) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := db.Where("id = ?", id).First(a.Bean).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a.Bean)
}

func (a *RestAPIs) Create(c *gin.Context) {
	if err := c.Bind(a.Bean); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// using reflect to set id field
	// u, _ := uuid.NewUUID()
	// rv := reflect.ValueOf(a.Bean).Elem()
	// rv.FieldByName("ID").SetString(u.String())

	if err := db.Create(a.Bean).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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

	// update all fields, no need for querying
	// err := db.Save(a.Bean).Error

	// update changed fields and we should query this record
	if err := db.Model(a.Bean).Updates(a.Bean).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.Where("id = ?", id).First(a.Bean).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a.Bean)
}

func (a *RestAPIs) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	rv := reflect.ValueOf(a.Bean).Elem()
	rv.FieldByName("ID").SetString(id)

	// soft delete
	// db.Delete(&bean)

	// delete permanently
	if err := db.Unscoped().Delete(a.Bean).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
