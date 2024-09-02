package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func InitWebServer() {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/v1/query/:func", func(c *gin.Context) {
		_func := c.Params.ByName("func")
		args := c.Request.URL.Query()["args"]
		data, err := Query(_func, args...)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    string(data),
		})
	})

	r.GET("/api/v1/invoke/:func", func(c *gin.Context) {
		_func := c.Params.ByName("func")
		args := c.Request.URL.Query()["args"]
		data, err := Query(_func, args...)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    string(data),
		})
	})

	r.GET("/api/v1/access/:id", connectToBackend)

	log.Info().Msg("Web server starting")
	go r.Run(":8080")
}
