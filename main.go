package main

import (
	"github.com/gin-gonic/gin"
	"san_francisco/config"
	"san_francisco/controller"
	"san_francisco/service"
	"san_francisco/utils"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	if config.Env == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(controller.CorsHandler())
	r.Use(controller.RequestLogger())
	r.Use(controller.APIKeyChecker())
	r.Use(controller.AuthChecker())
	r.Use(controller.ResponseLogger())
	r.Use(controller.SanFranciscoRoutes())
	return r
}

func main() {
	utils.InitializeLogger()
	defer utils.Logger.Sync()

	router = setupRouter()
	service.InitializeDB()
	service.RegisterRincon()
	service.GetAllAPIKeys()
	service.InitializeFirebase()
	service.ConnectDiscord()

	controller.InitializeRoutes(router)
	router.Run(":" + config.Port)
}
