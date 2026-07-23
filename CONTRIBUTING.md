# Contributing to go-ua-parser

Thank you for your interest in contributing! This guide will help you get started.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/go-ua-parser.git`
3. Create a branch: `git checkout -b feat/my-feature`
4. Make your changes
5. Run tests: `go test ./... -race -count=1`
6. Push and open a Pull Request

## Development

### Prerequisites

- Go 1.22+
- `gopkg.in/yaml.v3` (only if modifying rules and running the code generator)

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...                           # all tests
go test -run TestName ./...             # single test
go test -race ./...                     # race detector
go test -bench Benchmark -benchmem      # benchmarks
go test -fuzz FuzzParseString -fuzztime 30s  # fuzzing
```

### Code Generation

If you modify YAML rule files in `rules/`, regenerate the Go source:

```bash
cd cmd/uagen && go run .
```

This produces `rules_gen_*.go` files in the repo root. Commit both the YAML changes and regenerated Go files.

## Adding Rules

All rule files (`browsers.yaml`, `bots.yaml`, `devices.yaml`, `apps.yaml`, `engines.yaml`) share the same `match` field. Rules are matched against individual UA product tokens (name only), not the whole UA string; bot and device rules are additionally checked against the parenthesized-comment content of a token. An invalid `match: regex` pattern fails loudly rather than silently doing nothing: for builtin rules, `cmd/uagen` fatals at generation time naming the file and the offending token; for custom rules passed via `WithCustomBrowserRules`/`WithCustomBotRules`/`WithCustomDeviceRules`, `NewParser` returns a non-nil error naming the rule.

### Browsers

Edit `rules/browsers.yaml`:

```yaml
- token: "MyBrowser"        # UA token to match
  browser: "My Browser"     # display name
  family: "Chromium"        # browser family
  engine: "Blink"           # rendering engine
  match: exact              # exact | contains | prefix | regex (RE2)
```

Place specific rules (e.g., `EdgA`) **before** generic ones (e.g., `Chrome`). Order matters — first match wins.

### Bots

Edit `rules/bots.yaml`:

```yaml
- token: "MyBot"
  name: "My Bot"
  class: monitor             # search | social | ai | seo-tool | monitor | scraper | other
  vendor: "My Company"
  verified: false            # true for verified search engine bots
  match: exact
```

### Devices

Edit `rules/devices.yaml`:

```yaml
- token: "MyDevice"
  type: "mobile"             # mobile | tablet | desktop | console | smarttv | car | wearable
  vendor: "Acme"
  model: "Widget Pro"
  match: exact
```

### In-App Browsers

Edit `rules/apps.yaml`:

```yaml
- token: "MyApp"
  name: "My App"
  kind: "social"             # social | messaging | productivity | media | search | shopping | transport | finance
  match: exact
```

### Engines

Edit `rules/engines.yaml`:

```yaml
- token: "MyEngine"
  engine: "My Engine"
  match: exact
```

After editing any YAML file, run `cd cmd/uagen && go run .` and commit both files.

## Code Style

- Follow standard Go conventions (`go vet`, `gofmt`)
- Keep functions focused — one responsibility per function
- Add tests for new functionality
- Run benchmarks if touching the hot path to verify no regressions
- No external runtime dependencies in the main module

## Pull Request Guidelines

- Keep PRs focused on a single change
- Include tests for new features or bug fixes
- Update test fixtures if rule changes affect parsing behavior
- Run the full test suite before submitting
- Use conventional commit messages:
  - `feat:` for new features (triggers minor version bump)
  - `fix:` for bug fixes (triggers patch version bump)
  - `docs:` for documentation changes (no version bump)
  - `refactor:` for code changes that don't add features or fix bugs
  - `test:` for test-only changes
  - `chore:` for maintenance tasks

## Reporting Issues

- Use GitHub Issues
- Include the User-Agent string that caused the problem
- Include the expected vs actual parsing result
- If possible, include a minimal test case

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
