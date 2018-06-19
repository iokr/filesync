package router

import (
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/filesync/web/controller"
)


//Route 路由
func Route(router *gin.Engine) {
	//apiPrefix := config.ServerConfig.APIPrefix


	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	router.StaticFile("/favicon.ico", filepath.Join(getCurrentDirectory(), "./static/favicon.ico"))

	router.LoadHTMLGlob("views/**/*")

	api := router.Group("")
	{
		api.GET("/", controller.SigninGet)
		api.POST("/", controller.SigninPost)
	}

	adminAPI := router.Group("/admin")
	{
		adminAPI.GET("/index", controller.AdminIndexGet)

		adminAPI.GET("/task", controller.AdminTaskIndexGet)
		adminAPI.GET("/new_task", controller.AdminTaskCreateGet)
		adminAPI.POST("/new_task", controller.AdminTaskCreatePost)
		adminAPI.GET("/task/:id", controller.AdminTaskDisplay)
		adminAPI.GET("/task/:id/edit", controller.AdminTaskEditGet)
		adminAPI.POST("/task/:id/edit", controller.AdminTaskEditPost)
		adminAPI.POST("/task/:id/publish", controller.AdminTaskPublish)
		adminAPI.POST("/task/:id/delete", controller.AdminTaskDelete)

		adminAPI.GET("/tran_log/:page", controller.AdminTranLogGet)

		adminAPI.GET("/profile", controller.AdminProfileGet)
	}
}

func getCurrentDirectory() string {
	return ""
}


















