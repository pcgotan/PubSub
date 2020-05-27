package main

import (
	"github.com/gin-gonic/gin"
	
	"os"
	
)

func main() {
	//os.Setenv("PORT","3000")
	port := os.Getenv("PORT")
	
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	 router.Run(":" + port)
}