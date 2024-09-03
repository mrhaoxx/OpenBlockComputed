package main

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	_ "embed"
)

//go:embed setup.sh
var setupScript string

func InitWebServer() {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/api/v1/setup/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		c.Writer.Header().Set("Content-Type", "text/plain")

		c.String(200, strings.ReplaceAll(setupScript, "$1", id))
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
		data, err := Invoke(_func, args...)
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

	r.GET("/api/v1/listresources", func(c *gin.Context) {
		data, err := Query("ListComputeRes")
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
	r.GET("/api/v1/queryresource/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		data, err := Query("QueryComputeRes", id)
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

	r.GET("/api/v1/deleteresource/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		data, err := Invoke("DelComputeRes", id)
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

	r.GET("/api/v1/createresource/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		data, err := Invoke("CreateComputeRes", name)
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

	r.POST("/api/v1/updateresource/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")

		var result map[string]string

		if err := c.BindJSON(&result); err != nil {
			c.JSON(400, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return

		}

		_d, err := json.Marshal(result)

		if err != nil {
			c.JSON(400, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}

		data, err := InvokeTransistent("UpdateComputeRes", map[string][]byte{"update": _d}, id)

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

	r.POST("/api/v1/updateresourcessh/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")

		var result map[string]string

		if err := c.BindJSON(&result); err != nil {
			c.JSON(400, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return

		}

		_d, err := json.Marshal(result)

		if err != nil {
			c.JSON(400, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}

		data, err := InvokeTransistent("UpdateComputeRes", map[string][]byte{"ssh": _d}, id)

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

	r.NoRoute(func(c *gin.Context) {
		targetURL, err := url.Parse("https://10.196.109.185:8080")
		if err != nil {
			c.String(http.StatusInternalServerError, "Invalid target URL")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		c.Request.Host = targetURL.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	log.Info().Msg("Web server starting")
	go r.Run(":8080")
}
