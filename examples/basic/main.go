package main

import (
	"fmt"

	uax "github.com/motiv8-team/go-ua-parser"
)

func main() {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36"

	r := uax.Parse(ua)

	fmt.Printf("Browser: %s\n", r.ShortBrowser())
	fmt.Printf("Engine:  %s %s\n", r.Engine.Name, r.Engine.Version)
	fmt.Printf("OS:      %s\n", r.ShortOS())
	fmt.Printf("Device:  %s\n", r.DeviceClass())
	fmt.Printf("Is Bot:  %v\n", r.IsBot)
	fmt.Printf("Mobile:  %v\n", r.IsMobile)
}
