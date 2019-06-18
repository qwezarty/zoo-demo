package engine

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qwezarty/zoo-demo/models"
)

func Startup() *gorm.DB {
	dialect := "sqlite3"

	engine, err := gorm.Open(dialect, getConn(dialect))
	if err != nil {
		log.Fatalf("faltal error occour when conn to db: %v", err)
	}

	return engine.AutoMigrate(&models.Zoo{}, &models.Animal{})
}

func getConn(dialect string) string {
	if got := os.Getenv("DB"); got != "" {
		dialect = got
	}

	if dialect == "mssql" {
		return "xxxxxxxxxxxxxxxxxxxxxx"
	}

	return "./engine/zoo.db"
}
