package router

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jtieri/HabbGo/habbweb/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("client/", true))) // Enable static client files
	router.GET("/", controller.GetClient)

	return router
}
