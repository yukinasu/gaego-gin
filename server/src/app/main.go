package app

import (
	"gaego-gin/server/src/api"
	_ "gaego-gin/server/src/docs" // nolint
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title GAE/Go-Gin Sample API
// @version 1.0
// @description Sample API

// @license.name MIT

// @host localhost:8080
// @BasePath /api
func init() {
	r := gin.New()

	initAPI(r)
	initSwagger(r)

	http.Handle("/", r)
}

func initAPI(r *gin.Engine) {
	rg := r.Group("/api")
	api.SetupHoge(rg)
}

func initSwagger(r *gin.Engine) {
	rg := r.Group("/swagger")
	rg.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
