package router

import (
	"io"
	"log"
	"os"
	"templrjs/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Setup function
func Setup() *gin.Engine {
	r := gin.New()

	// Logging to a file.
	f, _ := os.Create("./logs/core.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	log.Println("Bootstrapping gin middlewares")
	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.GinContextToContextMiddleware())

	log.Println("Setting up routes")
	log.Println("Finished router setup")
	return r
}
