# go-ua-parser

A high-performance, production-grade HTTP User-Agent parser for Go with first-class bot and AI crawler detection.

[![Go Reference](https://pkg.go.dev/badge/github.com/motiv8-team/go-ua-parser.svg)](https://pkg.go.dev/github.com/motiv8-team/go-ua-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/motiv8-team/go-ua-parser)](https://goreportcard.com/report/github.com/motiv8-team/go-ua-parser)

## Features

- **Rich parsing**: Browser, Engine, OS, CPU, Device, In-App Browser, and Bot detection
- **170+ rules**: 39 browsers, 78 bots (search, social, AI, SEO, monitor), 26 devices, 22 in-app browsers
- **First-class bot detection**: Googlebot, GPTBot, ClaudeBot, PerplexityBot, AhrefsBot, social bots, monitoring tools, and 70+ more
- **Client Hints support**: Sec-CH-UA headers for accurate modern browser detection
- **High performance**: ~500K+ parses/sec per core, zero-alloc cache hits
- **Concurrency-safe**: Parser is immutable after construction, caches are sharded
- **Framework-ready**: Examples for Gin, Echo, Chi, Fiber + `ParseRequest(*http.Request)`
- **Extensible**: Custom rules, pre/post-parse hooks, browser channel detection
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
    fmt.Println(r.Engine.Name)     // "Blink"
    fmt.Println(r.OS.Name)         // "Windows"
    fmt.Println(r.DeviceClass())   // "desktop"
    fmt.Println(r.IsBot)           // false
}
```

## Installation

```bash
go get github.com/motiv8-team/go-ua-parser
```

Requires Go 1.22+.

## API Overview

### One-liner (global parser)

```go
r := uax.Parse("Mozilla/5.0 ... Chrome/123.0 ...")
```

### Reusable parser (recommended for servers)

```go
parser, _ := uax.NewParser()

r := parser.ParseString(ua)
r := parser.Parse(uax.Input{UAString: ua, ClientHints: ch})
r := parser.ParseRequest(httpReq) // extracts UA + Client Hints from *http.Request
```

### Bot detection

```go
bot := parser.DetectBot(ua)
if bot.IsBot {
    fmt.Println(bot.Name, bot.Class, bot.Vendor)
}
```

Bot classes: `search`, `social`, `ai`, `seo-tool`, `monitor`, `scraper`, `other`

### Caching

```go
// Single-mutex LRU (simple)
cached := uax.NewCachedParser(parser, 10000)

// Sharded LRU (high concurrency)
sharded := uax.NewShardedCache(parser, 16, 1000) // 16 shards, 1000 entries each
```

### Custom rules

```go
parser, _ := uax.NewParser(
    uax.WithCustomBotRules([]uax.BotRule{
        {Token: "InternalBot", Name: "Our Bot", Class: uax.BotMonitor, Vendor: "Us", Match: "exact"},
    }),
    uax.WithCustomBrowserRules([]uax.BrowserRule{
        {Token: "MyApp", Name: "My App", Family: "Chromium", Engine: "Blink", Match: "exact"},
    }),
)
```

### Hooks

```go
parser, _ := uax.NewParser(
    uax.WithPostParseHook(func(input uax.Input, result uax.Result, d time.Duration) {
        metrics.ParseDuration.Observe(d.Seconds())
        if result.IsBot {
            metrics.BotRequests.Inc()
        }
    }),
)
```

### HTTP middleware

```go
// net/http
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

See `examples/` for Gin, Echo, Chi, and Fiber middleware.

## Result Structure

```go
type Result struct {
    Browser   Browser  // Name, Version, Major, Family, Channel
    Engine    Engine   // Name, Version
    OS        OS       // Name, Version, Major, Minor, Patch
    CPU       CPU      // Architecture, Bits
    Device    Device   // Type, Vendor, Model, IsPhone, IsTablet, IsDesktop, IsTV
    App       App      // Name, Version, Kind (in-app browsers)
    Bot       Bot      // IsBot, Name, Class, Vendor, IsVerified, Confidence
    IsMobile  bool
    IsDesktop bool
    IsTablet  bool
    IsBot     bool
    IsCrawler bool
    IsInApp   bool
}
```

## Performance

Benchmarked on Apple M3 Pro:

| User-Agent | ns/op | allocs/op |
|------------|------:|----------:|
| curl (bot, shortest) | 201 | 0 |
| Googlebot | 1,121 | 2 |
| GPTBot (AI crawler) | 1,325 | 2 |
| Chrome Desktop | 1,852 | 5 |
| Edge Desktop | 1,947 | 5 |
| Firefox Desktop | 2,027 | 3 |
| Android Chrome | 2,150 | 5 |
| Safari Mobile | 2,529 | 7 |
| Samsung TV | 2,598 | 5 |
| Facebook In-App | 3,553 | 7 |

**Cache hit performance:**

| Cache Type | ns/op | allocs/op |
|------------|------:|----------:|
| LRU (single-mutex) | 32 | 0 |
| Sharded (16 shards) | 64 | 0 |
| Sharded (contended, parallel) | 82 | 0 |

Run benchmarks yourself:
```bash
go test -bench Benchmark -benchmem
```

## Code Generation

Rules are maintained in `rules/*.yaml` and compiled to Go source files:

```bash
go run cmd/uagen/*.go
```

Generated `rules_gen_*.go` files are committed to the repo. Library consumers never need to run the generator.

## Browser Channel Detection

Detects release channels from version patterns:

```go
r := parser.ParseString("Mozilla/5.0 ... Firefox/126.0a1")
r.Browser.Channel // "nightly"

r := parser.ParseString("Mozilla/5.0 ... Firefox/125.0b9")
r.Browser.Channel // "beta"
```

## Client Hints

When available, Client Hints provide more accurate detection than UA strings alone:

```go
r := parser.Parse(uax.Input{
    UAString: req.Header.Get("User-Agent"),
    ClientHints: uax.ClientHintsFromMap(map[string]string{
        "Sec-CH-UA":                  req.Header.Get("Sec-CH-UA"),
        "Sec-CH-UA-Platform":         req.Header.Get("Sec-CH-UA-Platform"),
        "Sec-CH-UA-Platform-Version": req.Header.Get("Sec-CH-UA-Platform-Version"),
    }),
})

// Or use ParseRequest which does this automatically:
r := parser.ParseRequest(req)
```

## License

MIT
