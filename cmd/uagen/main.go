package main

import (
	"fmt"
	"os"
	"regexp"
	"text/template"

	"gopkg.in/yaml.v3"
)

type browserRule struct {
	Token   string `yaml:"token"`
	Browser string `yaml:"browser"`
	Family  string `yaml:"family"`
	Engine  string `yaml:"engine"`
	Match   string `yaml:"match"`
}
type browserFile struct {
	Rules    []browserRule `yaml:"rules"`
	HasRegex bool          `yaml:"-"`
}

func (f *browserFile) prepare(source string) error {
	for _, r := range f.Rules {
		if err := checkRegex(source, r.Token, r.Match); err != nil {
			return err
		}
		if r.Match == "regex" {
			f.HasRegex = true
		}
	}
	return nil
}

type engineRule struct {
	Token  string `yaml:"token"`
	Engine string `yaml:"engine"`
	Match  string `yaml:"match"`
}
type engineFile struct {
	Rules    []engineRule `yaml:"rules"`
	HasRegex bool         `yaml:"-"`
}

func (f *engineFile) prepare(source string) error {
	for _, r := range f.Rules {
		if err := checkRegex(source, r.Token, r.Match); err != nil {
			return err
		}
		if r.Match == "regex" {
			f.HasRegex = true
		}
	}
	return nil
}

type botRule struct {
	Token    string `yaml:"token"`
	Name     string `yaml:"name"`
	Class    string `yaml:"class"`
	Vendor   string `yaml:"vendor"`
	Verified bool   `yaml:"verified"`
	Match    string `yaml:"match"`
}
type botFile struct {
	Rules    []botRule `yaml:"rules"`
	HasRegex bool      `yaml:"-"`
}

func (f *botFile) prepare(source string) error {
	for _, r := range f.Rules {
		if err := checkRegex(source, r.Token, r.Match); err != nil {
			return err
		}
		if r.Match == "regex" {
			f.HasRegex = true
		}
	}
	return nil
}

type deviceRule struct {
	Token  string `yaml:"token"`
	Type   string `yaml:"type"`
	Vendor string `yaml:"vendor"`
	Model  string `yaml:"model"`
	Match  string `yaml:"match"`
}
type deviceFile struct {
	Rules    []deviceRule `yaml:"rules"`
	HasRegex bool         `yaml:"-"`
}

func (f *deviceFile) prepare(source string) error {
	for _, r := range f.Rules {
		if err := checkRegex(source, r.Token, r.Match); err != nil {
			return err
		}
		if r.Match == "regex" {
			f.HasRegex = true
		}
	}
	return nil
}

type appRule struct {
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
	Kind  string `yaml:"kind"`
	Match string `yaml:"match"`
}
type appFile struct {
	Rules    []appRule `yaml:"rules"`
	HasRegex bool      `yaml:"-"`
}

func (f *appFile) prepare(source string) error {
	for _, r := range f.Rules {
		if err := checkRegex(source, r.Token, r.Match); err != nil {
			return err
		}
		if r.Match == "regex" {
			f.HasRegex = true
		}
	}
	return nil
}

// ruleFile is implemented by every *xxxFile YAML root struct. prepare is
// called once after unmarshalling: it validates that every `match: regex`
// rule's token compiles as RE2 (fatal-worthy otherwise) and records whether
// the file needs a "regexp" import in the generated output.
type ruleFile interface {
	prepare(source string) error
}

// checkRegex fails generation with a clear, file+token-qualified error when a
// rule has an unknown `match:` value or a `match: regex` token that is not
// valid RE2. Both checks run in prepare(), before os.Create truncates the
// output file, so a bad rule never leaves a partially-written rules_gen_*.go
// behind (the matchType funcMap re-checks during template execution as a
// backstop). Builtin regex rules are compiled once here, at generation time,
// so they arrive pre-compiled at package init in the generated rules_gen_*.go
// files — the matcher never compiles regexes lazily.
func checkRegex(source, token, match string) error {
	switch match {
	case "", "exact", "contains", "prefix":
		return nil
	case "regex":
		if _, err := regexp.Compile(token); err != nil {
			return fmt.Errorf("%s: rule %q: invalid regex: %w", source, token, err)
		}
		return nil
	default:
		return fmt.Errorf("%s: rule %q: invalid match type %q: must be exact, contains, prefix, or regex", source, token, match)
	}
}

var funcMap = template.FuncMap{
	"matchType": func(s string) (string, error) {
		switch s {
		case "", "exact":
			return "matchExact", nil
		case "contains":
			return "matchContains", nil
		case "prefix":
			return "matchPrefix", nil
		case "regex":
			return "matchRegex", nil
		default:
			return "", fmt.Errorf("invalid match type %q: must be exact, contains, prefix, or regex", s)
		}
	},
	"botClass": func(s string) string {
		switch s {
		case "search":
			return "BotSearch"
		case "social":
			return "BotSocial"
		case "monitor":
			return "BotMonitor"
		case "scraper":
			return "BotScraper"
		case "ai":
			return "BotAI"
		case "seo-tool":
			return "BotSEO"
		default:
			return "BotOther"
		}
	},
}

func generateFile(yamlPath, tmplStr, outPath string, data ruleFile) error {
	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", yamlPath, err)
	}
	if err := yaml.Unmarshal(raw, data); err != nil {
		return fmt.Errorf("parse %s: %w", yamlPath, err)
	}
	if err := data.prepare(yamlPath); err != nil {
		return err
	}
	tmpl, err := template.New("gen").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("generate %s from %s: %w", outPath, yamlPath, err)
	}
	return nil
}

func main() {
	fmt.Println("uagen: generating Go rule files from rules/*.yaml ...")

	var browsers browserFile
	if err := generateFile("../../rules/browsers.yaml", browserTmpl, "../../rules_gen_browser.go", &browsers); err != nil {
		fatal(err)
	}
	fmt.Printf("  browsers: %d rules -> rules_gen_browser.go\n", len(browsers.Rules))

	var engines engineFile
	if err := generateFile("../../rules/engines.yaml", engineTmpl, "../../rules_gen_engine.go", &engines); err != nil {
		fatal(err)
	}
	fmt.Printf("  engines:  %d rules -> rules_gen_engine.go\n", len(engines.Rules))

	var bots botFile
	if err := generateFile("../../rules/bots.yaml", botTmpl, "../../rules_gen_bot.go", &bots); err != nil {
		fatal(err)
	}
	fmt.Printf("  bots:     %d rules -> rules_gen_bot.go\n", len(bots.Rules))

	var devices deviceFile
	if err := generateFile("../../rules/devices.yaml", deviceTmpl, "../../rules_gen_device.go", &devices); err != nil {
		fatal(err)
	}
	fmt.Printf("  devices:  %d rules -> rules_gen_device.go\n", len(devices.Rules))

	var apps appFile
	if err := generateFile("../../rules/apps.yaml", appTmpl, "../../rules_gen_app.go", &apps); err != nil {
		fatal(err)
	}
	fmt.Printf("  apps:     %d rules -> rules_gen_app.go\n", len(apps.Rules))

	fmt.Println("uagen: done.")
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "uagen: %v\n", err)
	os.Exit(1)
}
