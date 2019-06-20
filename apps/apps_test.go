package apps

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qwezarty/zoo-demo/engine"
	"github.com/qwezarty/zoo-demo/models"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var apis RestAPIs
var w *httptest.ResponseRecorder
var r *http.Request

func setup() {
	engine.Sqlite = "../engine/zoo.db"
	db = engine.Startup("sqlite3", &models.Base{})
	Configure(db)

	w = httptest.NewRecorder()

	apis = RestAPIs{}
	apis.Bean = &models.Base{}
	apis.Beans = &[]models.Base{}

	router = gin.Default()
	router.GET("/base/:id", apis.Get)
	router.GET("/base", apis.List)
	router.POST("/base", apis.Create)
	router.PUT("/base/:id", apis.Update)
	router.DELETE("/base/:id", apis.Delete)
}

func teardown() {
	db.DropTable(&models.Base{})
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func mockCreate() string {
	u, _ := uuid.NewUUID()
	bean := &models.Base{ID: u.String()}
	db.Create(bean)

	return u.String()
}

func TestGet(t *testing.T) {
	bean := &models.Base{}

	id := mockCreate()
	r, _ = http.NewRequest("GET", "/base/"+id, nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	err := json.NewDecoder(w.Body).Decode(bean)
	assert.Nil(t, err)
	assert.Equal(t, id, bean.ID)

	db.Unscoped().Delete(bean)
}

func TestCreate(t *testing.T) {
	bean := &models.Base{}

	data, _ := json.Marshal(bean)
	r, _ = http.NewRequest("POST", "/base", bytes.NewBuffer(data))

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	err := json.NewDecoder(w.Body).Decode(bean)
	assert.Nil(t, err)
	assert.NotEmpty(t, bean.ID)

	db.Unscoped().Delete(bean)
}

func TestUpdate(t *testing.T) {
	id := mockCreate()
	bean := &models.Base{ID: id}

	data, _ := json.Marshal(bean)
	r, _ = http.NewRequest("PUT", "/base/"+id, bytes.NewBuffer(data))
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	err := json.NewDecoder(w.Body).Decode(bean)
	assert.Nil(t, err)
	assert.NotEmpty(t, bean.ID)
	if bean.UpdatedAt.IsZero() {
		t.Error("filed not updated")
	}

	db.Unscoped().Delete(bean)
}

func TestDelete(t *testing.T) {
	bean := &models.Base{}
	id := mockCreate()

	r, _ = http.NewRequest("DELETE", "/base/"+id, nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	db.Where("id = ?", id).First(bean)
	assert.Empty(t, bean.ID)
}

func TestList(t *testing.T) {
	records := []*models.Base{
		{ID: "0b58b99a-91a8-11e9-a956-1c872c7500f4"},
		{ID: "1f2bdc9f-925e-11e9-8705-1c872c7500f4"},
		{ID: "25525d58-925e-11e9-8705-1c872c7500f4"},
	}
	for _, rec := range records {
		db.Create(rec)
	}

	u, _ := url.Parse("/base")
	params := url.Values{}
	params.Add("page_token", "1")
	params.Add("page_size", "2")
	u.RawQuery = params.Encode()

	r, _ = http.NewRequest("GET", u.String(), nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	bean := &models.Pagination{}
	err := json.NewDecoder(w.Body).Decode(bean)
	assert.Nil(t, err)
	assert.Equal(t, 2, bean.PageSize)
	assert.Equal(t, 3, bean.TotalSize)
	assert.Equal(t, 1, bean.PageToken)
}
