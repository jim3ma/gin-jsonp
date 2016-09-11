package main

import (
	"github.com/gin-gonic/gin"
	jsonp "github.com/jim3ma/gin-jsonp"
)

func main() {
	r := gin.New()
	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(jsonp.Handler())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
		"message": "pong",
		})
	})
	r.Run(":8088") // listen and server on 0.0.0.0:8080
}
