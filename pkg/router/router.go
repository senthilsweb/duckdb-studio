package router

import (
	"io"
	"log"
	"os"
	"templrjs/pkg/config"
	"templrjs/pkg/duckdb"
	"templrjs/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Setup function
func Setup() *gin.Engine {
	r := gin.New()

	// Logging to a file.
	f, _ := os.Create("./logs/templrjs-core.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.Println("Application Name = [" + config.Config.Database.Dbname + "]")

	log.Println("Bootstrapping gin middlewares")
	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.GinContextToContextMiddleware())

	log.Println("Setting up routes")

	r.POST("/entities/:entity", duckdb.CreateEntity)
	r.GET("/entities/:entity", duckdb.GetEntities)
	r.PUT("/entities/:entity/:id", duckdb.UpdateEntity)
	r.DELETE("/entities/:entity/:id", duckdb.DeleteEntity)
	r.POST("/execute-query", duckdb.ExecuteCustomQuery)

	log.Println("Finished router setup")
	return r
}
