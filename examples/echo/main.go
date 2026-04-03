package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	uax "github.com/motiv8-team/go-ua-parser"
)

const uaKey = "ua_result"

func UAMiddleware(parser *uax.Parser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			result := parser.ParseString(c.Request().Header.Get("User-Agent"))
			c.Set(uaKey, result)
			return next(c)
		}
	}
}

func GetUA(c echo.Context) uax.Result {
	if v, ok := c.Get(uaKey).(uax.Result); ok {
		return v
	}
	return uax.Result{}
}

func main() {
	parser, _ := uax.NewParser()
	e := echo.New()
	e.Use(UAMiddleware(parser))

	e.GET("/", func(c echo.Context) error {
		ua := GetUA(c)
		return c.JSON(http.StatusOK, map[string]any{
			"browser": ua.ShortBrowser(),
			"os":      ua.ShortOS(),
			"device":  ua.DeviceClass(),
			"bot":     ua.IsBot,
		})
	})

	e.Start(":8080")
}
