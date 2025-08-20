package apiGateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gsq/music_bakcend_micorservice/middleware"
	"github.com/gsq/music_bakcend_micorservice/pkg/apiHandler"
)

func Setup() *gin.Engine {

	r := gin.Default()
	r.Static("/pages", "./static/pages")
	r.Static("/resources", "./static/resources")
	registerRoutes(r)

	return r
}

func registerRoutes(r *gin.Engine) {
	r.GET("/", apiHandler.Root)

	r.GET("/health", apiHandler.HealthCheck)

	api := r.Group("/api")
	{
		api.GET("/songs", apiHandler.GetSongs)
		api.GET("/songs/:name", apiHandler.GetSongbyName)
		api.POST("/login", apiHandler.Login)
		api.POST("/register", apiHandler.Register)
		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/check-auth", apiHandler.CheckAuth)
			auth.GET("/play-history", apiHandler.GetUserPlayHistory)
			auth.DELETE("/play-history/:songId", apiHandler.DeletePlayHistory)
			auth.POST("/songs/:songId/play", apiHandler.RecordPlayHistory)
			auth.POST("/logout", apiHandler.Logout)
		}
		api.GET("/songs/search", apiHandler.SearchSongs)
	}

}
