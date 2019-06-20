package engine

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qwezarty/zoo-demo/models"
)

var Sqlite = "./engine/zoo.db"

func Startup(dialect string, values ...interface{}) *gorm.DB {
	engine, err := gorm.Open(dialect, getConn(dialect))
	if err != nil {
		log.Fatalf("faltal error occour when conn to db: %v", err)
	}

	if len(values) == 0 {
		values = append(values, &models.Zoo{})
		values = append(values, &models.Animal{})
	}

	return engine.AutoMigrate(values...)
}

func getConn(dialect string) string {
	if got := os.Getenv("DB"); got != "" {
		dialect = got
	}

	switch dialect {
	case "mssql":
		return "conn string of mssql"
	case "sqlite3":
		return Sqlite
	default:
		return ""
	}
}
