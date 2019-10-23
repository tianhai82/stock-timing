package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Stock Timing starting")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://stock-timing.web.app"},
	}))

	if gin.Mode() == gin.DebugMode {
		r.Use(static.Serve("/", static.LocalFile("./web/public", false)))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		fmt.Println(err)
	}
}
