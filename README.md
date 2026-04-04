# go-ua-parser

A high-performance, production-grade HTTP User-Agent parser for Go with first-class bot and AI crawler detection.

[![Go Reference](https://pkg.go.dev/badge/github.com/motiv8-team/go-ua-parser.svg)](https://pkg.go.dev/github.com/motiv8-team/go-ua-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/motiv8-team/go-ua-parser)](https://goreportcard.com/report/github.com/motiv8-team/go-ua-parser)
[![CI](https://github.com/motiv8-team/go-ua-parser/actions/workflows/ci.yaml/badge.svg)](https://github.com/motiv8-team/go-ua-parser/actions/workflows/ci.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features

- **585 detection rules**: 124 browsers, 313 bots, 76 devices, 59 in-app browsers, 13 engines
- **35 OS types**: Windows (incl. Windows 11 via Client Hints), macOS, iOS, iPadOS, Android, ChromeOS, Linux distros (Ubuntu, Fedora, Debian, Arch, etc.), HarmonyOS, Tizen, KaiOS, FreeBSD, watchOS, tvOS, and more
- **First-class bot detection**: 313 bots across 7 classes — search engines, social bots, AI crawlers (GPTBot, ClaudeBot, DeepSeekBot, Bytespider), SEO tools, monitoring, scrapers, security scanners
- **14 heuristic patterns**: Catches unnamed bots via keywords (bot, crawler, spider, scraper, fetcher, archiver, validator, monitoring, etc.)
- **Client Hints support**: `Sec-CH-UA-*` headers for accurate modern browser detection, Windows 11 detection
- **Browser channel detection**: nightly, beta, dev, canary from version patterns
- **High performance**: ~500K+ parses/sec per core, zero-alloc cache hits at 32ns
- **Concurrency-safe**: Immutable Parser, sharded LRU cache (16 shards, 82ns contended)
- **Framework-ready**: Middleware examples for Gin, Echo, Chi, Fiber + `ParseRequest(*http.Request)`
- **Extensible**: Custom rules, pre/post-parse hooks
- **Pure Go**: Zero external runtime dependencies

## Quick Start

```go
package main

import (
    "fmt"
    uax "github.com/motiv8-team/go-ua-parser"
)

func main() {
    r := uax.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0.6312.86 Safari/537.36")

    fmt.Println(r.Browser.Name)    // "Chrome"
    fmt.Println(r.Browser.Version) // "123.0.6312.86"
    fmt.Println(r.Browser.Major)   // "123"
    fmt.Println(r.Browser.Family)  // "Chromium"
    fmt.Println(r.Engine.Name)     // "Blink"
    fmt.Println(r.OS.Name)         // "Windows"
    fmt.Println(r.CPU.Architecture)// "x86_64"
    fmt.Println(r.DeviceClass())   // "desktop"
    fmt.Println(r.IsBot)           // false
    fmt.Println(r.IsMobile)        // false
}
```

## Installation

```bash
go get github.com/motiv8-team/go-ua-parser
```

Requires Go 1.22+. Zero external dependencies.

## API

### Parse Functions

```go
// Global one-liner (uses sync.Once default parser)
r := uax.Parse(ua)

// Reusable parser (recommended for servers)
parser, err := uax.NewParser()
r := parser.ParseString(ua)

// With Client Hints
r := parser.Parse(uax.Input{UAString: ua, ClientHints: ch})

// From *http.Request (auto-extracts UA + all Client Hints headers)
r := parser.ParseRequest(httpReq)

// Zero-alloc variant (reuse caller-owned Result)
parser.ParseInto(input, &result)

// Bot-only fast path
bot := parser.DetectBot(ua)
```

### Result Structure

```go
type Result struct {
    Browser   Browser  // Name, Version, Major, Family, Channel
    Engine    Engine   // Name, Version
    OS        OS       // Name, Version, Major, Minor, Patch
    CPU       CPU      // Architecture, Bits
    Device    Device   // Type, Vendor, Model, IsPhone, IsTablet, IsDesktop, IsTV, IsTouch
    App       App      // Name, Version, Kind (in-app browsers)
    Bot       Bot      // IsBot, Name, Class, Vendor, Version, IsVerified, Confidence

    IsMobile  bool     // Device is a phone
    IsDesktop bool     // Device is desktop/laptop
    IsTablet  bool     // Device is a tablet
    IsBot     bool     // Any kind of bot/crawler
    IsCrawler bool     // Specifically a search/SEO crawler
    IsInApp   bool     // In-app browser (Facebook, Instagram, etc.)
}

// Convenience methods
r.ShortBrowser() // "Chrome 123"
r.ShortOS()      // "Windows 10"
r.DeviceClass()  // "mobile", "desktop", "tablet", "smarttv", "console", "car", "wearable", "unknown"
```

### Bot Detection

```go
bot := parser.DetectBot(ua)
if bot.IsBot {
    fmt.Println(bot.Name)       // "GPTBot"
    fmt.Println(bot.Class)      // "ai"
    fmt.Println(bot.Vendor)     // "OpenAI"
    fmt.Println(bot.IsVerified) // false
    fmt.Println(bot.Confidence) // 1.0
}
```

**Bot classes:** `search` | `social` | `ai` | `seo-tool` | `monitor` | `scraper` | `other`

**313 named bots** including:

| Class | Examples |
|-------|---------|
| Search | Googlebot, Bingbot, YandexBot, Baiduspider, DuckDuckBot, Applebot |
| Social | Twitterbot, facebookexternalhit, LinkedInBot, Slackbot, Discordbot, TelegramBot, WhatsApp |
| AI | GPTBot, ChatGPT-User, ClaudeBot, Claude-SearchBot, PerplexityBot, DeepSeekBot, Bytespider, CCBot, Amazonbot, Google-Extended, xAI-Bot |
| SEO | AhrefsBot, SemrushBot, MJ12bot, DotBot, Screaming Frog, SISTRIX, Seobility |
| Monitor | UptimeRobot, Pingdom, Datadog, NewRelic, Site24x7, GTmetrix, Uptime-Kuma, PRTG, Nagios |
| Scraper | curl, wget, python-requests, Scrapy, PhantomJS, HeadlessChrome, Playwright, Puppeteer, Shodan, Nmap |

Plus 14 heuristic keyword patterns that catch unnamed bots.

### Caching

```go
// Simple LRU cache (32ns cache hit, 0 allocs)
cached := uax.NewCachedParser(parser, 10000)
r := cached.ParseString(ua)

// Sharded cache for high concurrency (82ns contended, 0 allocs)
sharded := uax.NewShardedCache(parser, 16, 1000) // 16 shards x 1000 entries
r := sharded.ParseString(ua)

// Stats
stats := cached.Stats() // or sharded.Stats()
fmt.Println(stats.Hits, stats.Misses, stats.Size)
```

### Custom Rules

Override or extend builtin detection at parser creation time:

```go
parser, _ := uax.NewParser(
    uax.WithCustomBotRules([]uax.BotRule{
        {Token: "InternalBot", Name: "Our Bot", Class: uax.BotMonitor, Vendor: "Us", Match: "exact"},
    }),
    uax.WithCustomBrowserRules([]uax.BrowserRule{
        {Token: "MyApp", Name: "My App", Family: "Chromium", Engine: "Blink", Match: "exact"},
    }),
    uax.WithCustomDeviceRules([]uax.DeviceRule{
        {Token: "MyKiosk", Type: "embedded", Vendor: "Acme", Model: "Kiosk v2", Match: "exact"},
    }),
)
```

Custom rules are checked **before** builtins, so they can override default behavior.

### Hooks

```go
parser, _ := uax.NewParser(
    uax.WithPreParseHook(func(input uax.Input) {
        log.Println("Parsing:", input.UAString[:50])
    }),
    uax.WithPostParseHook(func(input uax.Input, result uax.Result, d time.Duration) {
        metrics.ParseDuration.Observe(d.Seconds())
        if result.IsBot {
            metrics.BotRequests.Inc()
        }
    }),
)
```

Hooks are nil-checked — zero cost when not configured.

### Client Hints

Modern Chromium browsers send reduced UA strings. Client Hints provide accurate detection:

```go
// Manual
r := parser.Parse(uax.Input{
    UAString: req.Header.Get("User-Agent"),
    ClientHints: uax.ClientHintsFromMap(map[string]string{
        "Sec-CH-UA":                  req.Header.Get("Sec-CH-UA"),
        "Sec-CH-UA-Platform":         req.Header.Get("Sec-CH-UA-Platform"),
        "Sec-CH-UA-Platform-Version": req.Header.Get("Sec-CH-UA-Platform-Version"),
        "Sec-CH-UA-Arch":             req.Header.Get("Sec-CH-UA-Arch"),
        "Sec-CH-UA-Model":            req.Header.Get("Sec-CH-UA-Model"),
        "Sec-CH-UA-Full-Version":     req.Header.Get("Sec-CH-UA-Full-Version"),
    }),
})

// Automatic (recommended)
r := parser.ParseRequest(req)
```

Client Hints fields override UA-derived values. Windows 11 is detected when `Sec-CH-UA-Platform-Version` >= 13.

### Browser Channel Detection

```go
r := parser.ParseString("Mozilla/5.0 ... Firefox/126.0a1")
r.Browser.Channel // "nightly"

r := parser.ParseString("Mozilla/5.0 ... Firefox/125.0b9")
r.Browser.Channel // "beta"
```

Detected channels: `nightly`, `beta`, `dev`, `canary` (from version patterns and Client Hints brands).

### Options

```go
parser, _ := uax.NewParser(
    uax.WithBotDetection(false),                    // disable bot detection
    uax.WithCustomBrowserRules([]uax.BrowserRule{}), // custom browser rules
    uax.WithCustomBotRules([]uax.BotRule{}),         // custom bot rules
    uax.WithCustomDeviceRules([]uax.DeviceRule{}),   // custom device rules
    uax.WithPreParseHook(fn),                        // pre-parse callback
    uax.WithPostParseHook(fn),                       // post-parse callback with timing
)
```

## Framework Middleware

Each framework example is a separate Go module (no framework deps in the main library).

### net/http (stdlib)

```go
func UAMiddleware(parser *uax.Parser) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            result := parser.ParseRequest(r)
            ctx := context.WithValue(r.Context(), uaKey, result)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Gin

```go
func UAMiddleware(parser *uax.Parser) gin.HandlerFunc {
    return func(c *gin.Context) {
        result := parser.ParseString(c.GetHeader("User-Agent"))
        c.Set("ua_result", result)
        c.Next()
    }
}
```

### Echo

```go
func UAMiddleware(parser *uax.Parser) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            result := parser.ParseString(c.Request().Header.Get("User-Agent"))
            c.Set("ua_result", result)
            return next(c)
        }
    }
}
```

### Chi

```go
func UAMiddleware(parser *uax.Parser) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            result := parser.ParseRequest(r)
            ctx := context.WithValue(r.Context(), ctxKey{}, result)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Fiber

```go
func UAMiddleware(parser *uax.Parser) fiber.Handler {
    return func(c *fiber.Ctx) error {
        result := parser.ParseString(c.Get("User-Agent"))
        c.Locals("ua_result", result)
        return c.Next()
    }
}
```

Complete working examples in `examples/gin/`, `examples/echo/`, `examples/chi/`, `examples/fiber/`.

## Performance

Benchmarked on Apple M3 Pro (arm64):

**Parse performance:**

| User-Agent Type | ns/op | B/op | allocs/op |
|----------------|------:|-----:|----------:|
| curl (short bot) | 284 | 0 | 0 |
| Googlebot | 1,121 | 112 | 2 |
| GPTBot (AI) | 1,325 | 144 | 2 |
| Chrome Desktop | 1,852 | 272 | 5 |
| Edge Desktop | 1,947 | 288 | 5 |
| Firefox Desktop | 2,027 | 160 | 3 |
| Android Chrome | 2,150 | 272 | 5 |
| Safari Mobile | 2,529 | 336 | 7 |
| Samsung TV | 2,598 | 264 | 5 |
| Facebook In-App | 3,553 | 336 | 7 |

**Cache performance:**

| Cache Type | ns/op | allocs/op |
|------------|------:|----------:|
| LRU cache hit | 32 | 0 |
| Sharded cache hit | 64 | 0 |
| Sharded contended (parallel) | 82 | 0 |

**Internal components:**

| Component | ns/op | allocs/op |
|-----------|------:|----------:|
| Tokenizer | 75 | 0 |
| Trie lookup | 4 | 0 |

Run benchmarks:
```bash
go test -bench Benchmark -benchmem -count=3
```

## Architecture

Three-stage parsing pipeline:

```
UA String ──→ [Tokenizer] ──→ [Matcher] ──→ [Assembler] ──→ Result
                  │                │              │
            zero-alloc      trie (4ns)     Client Hints
            fixed buffer    + regex        merge + derive
```

1. **Tokenize** (`tokenizer.go`): Zero-alloc scanner extracts product tokens (name/version) and comment blocks using a fixed `[24]token` buffer.
2. **Match** (`match_*.go`, `matcher.go`, `trie.go`): Hybrid trie + linear-scan matcher. Trie handles exact-match rules in O(key-length); contains/prefix rules fall back to indexed scan.
3. **Assemble** (`merge.go`): Merges Client Hints (CH takes precedence), computes convenience booleans.

Rules are defined in `rules/*.yaml` and compiled to Go source via `cmd/uagen`. Generated `rules_gen_*.go` files are committed — consumers never run the generator.

## Coverage

| Category | Count | Description |
|----------|------:|-------------|
| Browsers | 124 | Chrome, Firefox, Safari, Edge, Opera, Brave, Vivaldi, Samsung Internet, Arc, Whale, DuckDuckGo, 100+ more |
| Engines | 13 | Blink, WebKit, Gecko, Trident, Presto, EdgeHTML, Goanna, KHTML, Servo, NetSurf, and more |
| Bots | 313 | Search, social, AI, SEO, monitor, scraper, security scanners, HTTP libraries |
| Devices | 76 | iPhones, iPads, Samsung Galaxy, Pixel, Xiaomi, consoles (PS5, Xbox, Switch), TVs, Kindle, Tesla, Apple Watch |
| In-App | 59 | Facebook, Instagram, TikTok, WeChat, Telegram, Slack, Teams, Spotify, and 50+ more |
| OS | 35 | Windows (10/11), macOS, iOS, iPadOS, Android, ChromeOS, 15+ Linux distros, HarmonyOS, Tizen, KaiOS, BSD variants |
| Heuristics | 14 | bot, crawler, spider, scraper, fetcher, archiver, validator, monitoring, analyzer, and more |

## Code Generation

Rules are maintained in YAML and compiled to Go:

```bash
cd cmd/uagen && go run .
```

This reads `rules/*.yaml` and generates `rules_gen_*.go` in the repo root. Generated files are committed — library consumers never run the generator or need `gopkg.in/yaml.v3`.

To add a new rule, edit the appropriate YAML file and regenerate:

```yaml
# rules/bots.yaml
- token: "MyNewBot"
  name: "My New Bot"
  class: monitor
  vendor: "My Company"
  match: exact
```

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make changes, add tests
4. Run `go test ./... -race -count=1`
5. If you changed YAML rules: `cd cmd/uagen && go run .`
6. Commit and push
7. Open a Pull Request

### Adding Rules

- **Browser**: Edit `rules/browsers.yaml` — token, browser name, family, engine, match type
- **Bot**: Edit `rules/bots.yaml` — token, name, class (search/social/ai/seo-tool/monitor/scraper/other), vendor, match type
- **Device**: Edit `rules/devices.yaml` — token, type (mobile/tablet/console/smarttv/car/wearable), vendor, model, match type
- **In-App**: Edit `rules/apps.yaml` — token, app name, kind, match type
- **Engine**: Edit `rules/engines.yaml` — token, engine name, match type

Then run the generator and commit both the YAML and generated Go files.

## License

[MIT](LICENSE)
