package main

import (
	"github.com/gofiber/fiber/v2"
	uax "github.com/motiv8-team/go-ua-parser"
)

const uaKey = "ua_result"

func UAMiddleware(parser *uax.Parser) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result := parser.ParseString(c.Get("User-Agent"))
		c.Locals(uaKey, result)
		return c.Next()
	}
}

func GetUA(c *fiber.Ctx) uax.Result {
	if v, ok := c.Locals(uaKey).(uax.Result); ok {
		return v
	}
	return uax.Result{}
}

func main() {
	parser, _ := uax.NewParser()
	app := fiber.New()
	app.Use(UAMiddleware(parser))

	app.Get("/", func(c *fiber.Ctx) error {
		ua := GetUA(c)
		return c.JSON(fiber.Map{
			"browser": ua.ShortBrowser(),
			"os":      ua.ShortOS(),
			"device":  ua.DeviceClass(),
			"bot":     ua.IsBot,
		})
	})

	app.Listen(":8080")
}
