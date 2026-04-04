// Package uaxotel provides OpenTelemetry integration helpers for go-ua-parser.
package uaxotel

import (
	uax "github.com/motiv8-team/go-ua-parser"
	"go.opentelemetry.io/otel/attribute"
)

const (
	attrBrowserName    = attribute.Key("http.user_agent.browser.name")
	attrBrowserVersion = attribute.Key("http.user_agent.browser.version")
	attrBrowserFamily  = attribute.Key("http.user_agent.browser.family")
	attrEngineName     = attribute.Key("http.user_agent.engine.name")
	attrEngineVersion  = attribute.Key("http.user_agent.engine.version")
	attrOSName         = attribute.Key("http.user_agent.os.name")
	attrOSVersion      = attribute.Key("http.user_agent.os.version")
	attrDeviceType     = attribute.Key("http.user_agent.device.type")
	attrDeviceVendor   = attribute.Key("http.user_agent.device.vendor")
	attrDeviceModel    = attribute.Key("http.user_agent.device.model")
	attrCPUArch        = attribute.Key("http.user_agent.cpu.architecture")
	attrBotName        = attribute.Key("http.user_agent.bot.name")
	attrBotClass       = attribute.Key("http.user_agent.bot.class")
	attrIsBot          = attribute.Key("http.user_agent.is_bot")
	attrIsMobile       = attribute.Key("http.user_agent.is_mobile")
)

// Attributes returns OpenTelemetry span attributes for the parsed result.
// Only non-empty fields are included to keep spans compact.
func Attributes(r uax.Result) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, 16)

	if r.Browser.Name != "" {
		attrs = append(attrs, attrBrowserName.String(r.Browser.Name))
	}
	if r.Browser.Version != "" {
		attrs = append(attrs, attrBrowserVersion.String(r.Browser.Version))
	}
	if r.Browser.Family != "" {
		attrs = append(attrs, attrBrowserFamily.String(r.Browser.Family))
	}
	if r.Engine.Name != "" {
		attrs = append(attrs, attrEngineName.String(r.Engine.Name))
	}
	if r.Engine.Version != "" {
		attrs = append(attrs, attrEngineVersion.String(r.Engine.Version))
	}
	if r.OS.Name != "" {
		attrs = append(attrs, attrOSName.String(r.OS.Name))
	}
	if r.OS.Version != "" {
		attrs = append(attrs, attrOSVersion.String(r.OS.Version))
	}
	if r.Device.Type != "" {
		attrs = append(attrs, attrDeviceType.String(r.Device.Type))
	}
	if r.Device.Vendor != "" {
		attrs = append(attrs, attrDeviceVendor.String(r.Device.Vendor))
	}
	if r.Device.Model != "" {
		attrs = append(attrs, attrDeviceModel.String(r.Device.Model))
	}
	if r.CPU.Architecture != "" {
		attrs = append(attrs, attrCPUArch.String(r.CPU.Architecture))
	}
	if r.IsBot {
		attrs = append(attrs, attrIsBot.Bool(true))
		if r.Bot.Name != "" {
			attrs = append(attrs, attrBotName.String(r.Bot.Name))
		}
		if r.Bot.Class != "" {
			attrs = append(attrs, attrBotClass.String(string(r.Bot.Class)))
		}
	}
	if r.IsMobile {
		attrs = append(attrs, attrIsMobile.Bool(true))
	}

	return attrs
}
