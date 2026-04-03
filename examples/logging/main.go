package main

import (
	"fmt"
	"log/slog"
	"os"

	uax "github.com/motiv8-team/go-ua-parser"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	uas := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0 Safari/537.36",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"curl/8.4.0",
	}

	parser, _ := uax.NewParser()
	for _, ua := range uas {
		r := parser.ParseString(ua)
		logger.Info("request",
			slog.String("browser", r.ShortBrowser()),
			slog.String("os", r.ShortOS()),
			slog.String("device", r.DeviceClass()),
			slog.Bool("bot", r.IsBot),
			slog.String("bot_name", r.Bot.Name),
			slog.String("bot_class", string(r.Bot.Class)),
		)
	}

	fmt.Println("\n(Above is structured JSON log output)")
}
