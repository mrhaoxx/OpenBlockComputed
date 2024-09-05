package main

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	_ "embed"
)

//go:embed setup.sh
var setupScript string

func InitWebServer() {

	port := "8080"

	TestEnv("PORT", &port)

	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/api/v1/setup/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		c.Writer.Header().Set("Content-Type", "text/plain")

		_r := strings.ReplaceAll(setupScript, "$1", id)
		_r = strings.ReplaceAll(_r, "$2", c.Request.Host)

		c.String(200, _r)
	})

	r.GET("/api/v1/query/:func", func(c *gin.Context) {
		_func := c.Params.ByName("func")
		args := c.Request.URL.Query()["args"]
		data, err := Query(_func, args...)
		if err != nil {
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
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
	r.GET("/api/v1/accesslogs/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		data, err := Query("GetConnectionLogs", id)
		if err != nil {
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/whoami", func(c *gin.Context) {
		data, err := Query("GetUserInfo")
		if err != nil {
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return

		}

		_d, err := json.Marshal(result)

		if err != nil {
			c.JSON(200, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}

		data, err := InvokeTransistent("UpdateComputeRes", map[string][]byte{"update": _d}, id)

		if err != nil {
			c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return

		}

		_d, err := json.Marshal(result)

		if err != nil {
			c.JSON(200, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}

		data, err := InvokeTransistent("UpdateComputeRes", map[string][]byte{"ssh": _d}, id)

		if err != nil {
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/claimresource/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		data, err := Invoke("ClaimRent", id)
		if err != nil {
			c.JSON(200, gin.H{
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

	r.POST("/api/v1/market/put/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")

		var result map[string]string

		if err := c.BindJSON(&result); err != nil {
			c.JSON(200, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}

		t := result["duration"]

		if t == "" {
			t = "0"
		}

		_t, err := time.ParseDuration(t)

		if err != nil {
			c.JSON(200, gin.H{
				"message": "error",
				"error":   err.Error(),
			})
			return
		}

		data, err := Invoke("PutOnMarket", id, strconv.Itoa(int(_t.Microseconds())), result["price"])

		if err != nil {
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/market/get/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")

		data, err := Query("GetMarketElement", id)

		if err != nil {
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/market/price/:id/:price", func(c *gin.Context) {
		id := c.Params.ByName("id")
		price := c.Params.ByName("price")

		data, err := Invoke("MakePrice", id, price)

		if err != nil {
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/market/lock/:id/:winner/:price", func(c *gin.Context) {
		id := c.Params.ByName("id")
		winner := c.Params.ByName("winner")
		price := c.Params.ByName("price")

		data, err := Invoke("LockMarketElement", id, winner, price)

		if err != nil {
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/market/end/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")

		data, err := Invoke("EndMarketElement", id)

		if err != nil {
			c.JSON(200, gin.H{
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

	r.GET("/api/v1/market/list", func(c *gin.Context) {

		data, err := Query("ListMarketElements")

		if err != nil {
			c.JSON(200, gin.H{
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
	go r.Run(":" + port)
}
