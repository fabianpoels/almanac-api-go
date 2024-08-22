package server

import (
	"almanac-api/controllers"
	"almanac-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// FOR DEV PURPOSES!!!!
	// TODO: rework to dynamically configure cors, depending on the env
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:9000", "http://127.0.0.1:9000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	auth := new(controllers.AuthController)
	// user := new(controllers.UserController)

	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.POST("/auth/login", auth.Login)
			v1.POST("/auth/refresh-token", auth.RefreshToken)
			v1.Use(middleware.ValidateJwt())
			{
				// v1.GET("workspaces", user.Workspaces)
			}
		}
	}
	return router

}
