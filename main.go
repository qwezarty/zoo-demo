package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/qwezarty/zoo-demo/apps"
	"github.com/qwezarty/zoo-demo/apps/animal"
	"github.com/qwezarty/zoo-demo/apps/zoo"
	"github.com/qwezarty/zoo-demo/engine"
)

func main() {
	// startup all managers
	var router = gin.Default()
	var db = engine.Startup("sqlite3")

	// register all sub-routes
	apps.Configure(db)
	zoo.Configure(router, db)
	animal.Configure(router, db)

	log.Fatal(router.Run("0.0.0.0:30096"))
}
