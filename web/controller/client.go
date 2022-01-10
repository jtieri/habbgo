package controller

import "github.com/gin-gonic/gin"

func GetClient(c *gin.Context) {
	c.File("./client/client.html")
}
