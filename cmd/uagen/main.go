package main

import (
	"fmt"
	"os"
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
type browserFile struct{ Rules []browserRule `yaml:"rules"` }

type engineRule struct {
	Token  string `yaml:"token"`
	Engine string `yaml:"engine"`
	Match  string `yaml:"match"`
}
type engineFile struct{ Rules []engineRule `yaml:"rules"` }

type botRule struct {
	Token    string `yaml:"token"`
	Name     string `yaml:"name"`
	Class    string `yaml:"class"`
	Vendor   string `yaml:"vendor"`
	Verified bool   `yaml:"verified"`
	Match    string `yaml:"match"`
}
type botFile struct{ Rules []botRule `yaml:"rules"` }

type deviceRule struct {
	Token  string `yaml:"token"`
	Type   string `yaml:"type"`
	Vendor string `yaml:"vendor"`
	Model  string `yaml:"model"`
	Match  string `yaml:"match"`
}
type deviceFile struct{ Rules []deviceRule `yaml:"rules"` }

type appRule struct {
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
	Kind  string `yaml:"kind"`
	Match string `yaml:"match"`
}
type appFile struct{ Rules []appRule `yaml:"rules"` }

var funcMap = template.FuncMap{
	"matchType": func(s string) string {
		switch s {
		case "exact":
			return "matchExact"
		case "contains":
			return "matchContains"
		case "prefix":
			return "matchPrefix"
		case "regex":
			return "matchRegex"
		default:
			return "matchExact"
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

func generateFile[T any](yamlPath, tmplStr, outPath string, data *T) error {
	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", yamlPath, err)
	}
	if err := yaml.Unmarshal(raw, data); err != nil {
		return fmt.Errorf("parse %s: %w", yamlPath, err)
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
	return tmpl.Execute(f, data)
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
