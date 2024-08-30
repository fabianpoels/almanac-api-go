package server

import (
	"almanac-api/config"
	"almanac-api/controllers"
	"almanac-api/controllers/admin"
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
	corsConfig := cors.DefaultConfig()

	if config.GetEnv("ENVIRONMENT") == "dev" {
		// LOCAL DEV CONFIG
		router.SetTrustedProxies(nil)
		corsConfig.AllowOrigins = []string{"http://localhost:9000", "http://127.0.0.1:9000"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
		corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
	} else {
		domain := config.GetEnv("DOMAIN")
		router.SetTrustedProxies(nil)
		corsConfig.AllowOrigins = []string{domain}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
		corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
	}

	router.Use(cors.New(corsConfig))

	auth := new(controllers.AuthController)
	public := new(controllers.PublicController)

	// ADMIN CONTROLLERS
	newsItems := new(admin.NewsItemsController)

	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			// public routes
			v1.GET("/news", public.NewsItems)
			v1.GET("/categories", public.Categories)

			// auth
			v1.POST("/auth/login", auth.Login)
			v1.POST("/auth/refresh-token", auth.RefreshToken)
			v1.POST("/auth/logout", middleware.ValidateJwt(), auth.Logout)
			v1.Use(middleware.ValidateJwt())
			{
				// LOGGED IN ROUTES

				// ADMIN ROUTES
				admin := v1.Group("a")
				{
					admin.Use(middleware.ValidateAdmin())
					{
						news := admin.Group("news")
						{
							news.GET("", newsItems.List)
							// news.PATCH("/newsItem", newsItems.Patch)

						}
					}
				}
			}
		}
	}
	return router

}
