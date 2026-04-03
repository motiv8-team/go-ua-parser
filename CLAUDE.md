# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`go-ua-parser` — a high-performance, zero-allocation HTTP User-Agent parser library for Go.
Parses UA strings and Client Hints into rich structured results: browser, engine, OS, CPU, device, in-app browser, and bot/crawler/AI-agent detection.

## Build & Test

```bash
go build ./...
go test ./...
go test -run TestName ./...          # single test
go test -bench Benchmark -benchmem   # benchmarks
go test -fuzz FuzzParseString -fuzztime 30s  # fuzz testing
go vet ./...
```

## Code Generation

Rules are defined in `rules/*.yaml` and compiled to Go source via:
```bash
cd cmd/uagen && go run .
```
Or from the repo root: `go generate ./...` (uses the `//go:generate` directive in `uax.go`).

`cmd/uagen` is a separate Go module with its own `go.mod` (to keep `gopkg.in/yaml.v3` out of the main module). It reads all YAML files in `rules/` and generates `rules_gen_*.go` files. Generated files are committed to the repo — library consumers never run the generator.

Current rule counts: ~40 browsers, ~80 bots, ~26 devices, ~22 apps, 5 engines.

## Architecture

Three-stage parsing pipeline:
1. **Tokenize** (`tokenizer.go`): Zero-alloc scanner splits UA into product tokens (name/version) and comment blocks. Uses a fixed-size buffer of 24 tokens.
2. **Match** (`match_*.go`, `matcher.go`, `trie.go`): Hybrid trie + regex rule matcher. Each category has its own file: `match_browser.go`, `match_engine.go`, `match_os.go`, `match_device.go`, `match_cpu.go`, `match_bot.go`, `match_app.go`.
3. **Assemble** (`merge.go`): Merges Client Hints (CH takes precedence), computes convenience booleans (`IsMobile`, `IsBot`, etc.).

Key types:
- `Parser` — immutable after construction, safe for concurrent use
- `Result` — full parse output with `Browser`, `Engine`, `OS`, `CPU`, `Device`, `App`, `Bot`
- `CachedParser` — single-mutex LRU cache wrapper
- `ShardedCache` — 16-shard concurrent cache for high-throughput servers
- `Input` — holds UA string + optional `ClientHints`

Additional entry points:
- `ParseRequest(*http.Request)` — extracts UA + Client Hints headers automatically
- `DetectBot(ua string)` — bot-only fast path

## Extensibility

- **Custom rules**: `WithCustomBrowserRules()`, `WithCustomBotRules()`, `WithCustomDeviceRules()` — prepend user rules before builtins
- **Hooks**: `WithPreParseHook()`, `WithPostParseHook()` — callbacks for instrumentation/metrics
- **Browser channels**: `Browser.Channel` detects nightly/beta/dev/canary from version patterns

## Design Principles

- Zero heap allocations on hot path after parser initialization
- Rules: trie for fast exact match, regex only for long-tail/complex patterns
- `Result` is a value type returned by copy (stack-allocated when possible)
- `ParseInto(*Result)` variant for callers who want to reuse a Result
- Bot detection is first-class: `BotClass` enum covers search, social, monitor, scraper, ai, seo-tool
- Client Hints override UA-derived fields (browser name/version, platform, arch, model)
- Generated rule files (`rules_gen_*.go`) come from YAML via `cmd/uagen` — edit YAML, not generated Go

## Framework Examples

Separate Go modules in `examples/` (don't pull framework deps into main lib):
- `examples/gin/` — Gin middleware
- `examples/echo/` — Echo middleware
- `examples/chi/` — Chi middleware
- `examples/fiber/` — Fiber middleware
- `examples/basic/` — stdlib one-liner
- `examples/middleware/` — stdlib net/http middleware
- `examples/logging/` — structured logging with slog
