package controller

import "github.com/gin-gonic/gin"

func GetClient(c *gin.Context) {
	//fmt.Println("Inside GetIndex() .....")
	c.File("./client/client.html")
}
