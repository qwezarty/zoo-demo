package apps

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

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
	defer db.Unscoped().Delete(bean)

	r, _ = http.NewRequest("GET", "/base/"+id, nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	err := json.NewDecoder(w.Body).Decode(bean)
	assert.Nil(t, err)
	assert.Equal(t, id, bean.ID)
}

func TestCreate(t *testing.T) {
	bean := &models.Base{}
	defer db.Unscoped().Delete(bean)

	data, _ := json.Marshal(bean)
	r, _ = http.NewRequest("POST", "/base", bytes.NewBuffer(data))

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	err := json.NewDecoder(w.Body).Decode(bean)
	assert.Nil(t, err)
	assert.NotEmpty(t, bean.ID)
}

func TestUpdate(t *testing.T) {
	id := mockCreate()
	bean := &models.Base{ID: id}
	defer db.Unscoped().Delete(bean)

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
}

func TestDelete(t *testing.T) {
	bean := &models.Base{}
	id := mockCreate()
	defer db.Unscoped().Delete(bean)

	r, _ = http.NewRequest("DELETE", "/base/"+id, nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	db.Where("id = ?", id).First(bean)
	assert.Empty(t, bean.ID)
}

func TestList(t *testing.T) {
	now := time.Now()
	records := []*models.Base{
		{ID: "8c5b25ae-91a8-11e9-a956-1c872c7500f4", CreatedAt: now.Add(-2 * time.Hour)},
		{ID: "0b58b99a-91a8-11e9-a956-1c872c7500f4", CreatedAt: now.Add(-1 * time.Hour)},
		{ID: "1f2bdc9f-925e-11e9-8705-1c872c7500f4", CreatedAt: now.Add(1 * time.Hour)},
		{ID: "25525d58-925e-11e9-8705-1c872c7500f4", CreatedAt: now.Add(2 * time.Hour)},
	}
	for _, rec := range records {
		db.Create(rec)
	}
	defer func() {
		for _, rec := range records {
			db.Unscoped().Delete(rec)
		}
	}()

	u0, _ := url.Parse("/base")

	u1, _ := url.Parse("/base")
	params1 := url.Values{}
	params1.Add("page_token", "1")
	params1.Add("page_size", "3")
	u1.RawQuery = params1.Encode()

	u2, _ := url.Parse("/base")
	params2 := url.Values{}
	params2.Add("start_time", now.Format(time.RFC3339))
	u2.RawQuery = params2.Encode()

	testCases := []struct {
		url   string
		count int
		total int
	}{
		{url: u0.String(), count: 4, total: 4},
		{url: u1.String(), count: 1, total: 4},
		{url: u2.String(), count: 2, total: 2},
	}

	for _, tc := range testCases {
		r, _ = http.NewRequest("GET", tc.url, nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code)

		bean := &models.Pagination{}
		err := json.NewDecoder(w.Body).Decode(bean)
		assert.Nil(t, err)
		assert.Equal(t, tc.total, bean.TotalSize)

		entities := bean.Body.([]interface{})
		assert.Equal(t, tc.count, len(entities))
	}
}
