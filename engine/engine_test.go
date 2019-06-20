package engine

import (
	"testing"

	"github.com/qwezarty/zoo-demo/models"
	"github.com/stretchr/testify/assert"
)

func TestSqlite(t *testing.T) {
	Sqlite = "./zoo.db"

	bean := &models.Base{}
	db := Startup("sqlite3", bean)
	db.DropTable(bean)

	Sqlite = "./engine/zoo.db"
}

func TestGetConn(t *testing.T) {
	testCases := []struct {
		dialect string
		want    string
	}{
		{dialect: "mssql", want: "conn string of mssql"},
		{dialect: "sqlite3", want: "./engine/zoo.db"},
		{dialect: "mysql", want: ""},
	}

	for _, tc := range testCases {
		got := getConn(tc.dialect)
		assert.Equal(t, tc.want, got)
	}
}
