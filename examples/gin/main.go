package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	uax "github.com/motiv8-team/go-ua-parser"
)

const uaKey = "ua_result"

func UAMiddleware(parser *uax.Parser) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := parser.ParseString(c.GetHeader("User-Agent"))
		c.Set(uaKey, result)
		c.Next()
	}
}

func GetUA(c *gin.Context) uax.Result {
	if v, ok := c.Get(uaKey); ok {
		return v.(uax.Result)
	}
	return uax.Result{}
}

func main() {
	parser, _ := uax.NewParser()
	r := gin.Default()
	r.Use(UAMiddleware(parser))

	r.GET("/", func(c *gin.Context) {
		ua := GetUA(c)
		c.JSON(http.StatusOK, gin.H{
			"browser": ua.ShortBrowser(),
			"os":      ua.ShortOS(),
			"device":  ua.DeviceClass(),
			"bot":     ua.IsBot,
		})
	})

	r.Run(":8080")
}
