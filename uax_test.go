package uax

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestResultZeroValue(t *testing.T) {
	var r Result
	if r.IsMobile {
		t.Error("zero Result should not be mobile")
	}
	if r.IsBot {
		t.Error("zero Result should not be bot")
	}
	if r.Browser.Name != "" {
		t.Error("zero Result browser name should be empty")
	}
	if r.Device.Type != "" {
		t.Error("zero Result device type should be empty")
	}
}

func TestResultJSON(t *testing.T) {
	r := Result{
		Browser: Browser{Name: "Chrome", Version: "123.0.6312.86", Major: "123"},
		OS:      OS{Name: "Windows", Version: "10.0", Major: "10"},
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var r2 Result
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r2.Browser.Name != "Chrome" {
		t.Errorf("got browser %q, want Chrome", r2.Browser.Name)
	}
	if r2.OS.Major != "10" {
		t.Errorf("got OS major %q, want 10", r2.OS.Major)
	}
}

func TestNewParserDefaults(t *testing.T) {
	p, err := NewParser()
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParserParseStringEmpty(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("")
	if r.Browser.Name != "" {
		t.Errorf("empty UA should produce empty browser, got %q", r.Browser.Name)
	}
	if r.IsBot {
		t.Error("empty UA should not be bot")
	}
}

func TestParserParseStringChrome(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36")
	if r.Browser.Name != "Chrome" {
		t.Errorf("browser = %q, want Chrome", r.Browser.Name)
	}
	if r.Browser.Major != "123" {
		t.Errorf("major = %q, want 123", r.Browser.Major)
	}
	if r.Engine.Name != "Blink" {
		t.Errorf("engine = %q, want Blink", r.Engine.Name)
	}
	if r.OS.Name != "Windows" {
		t.Errorf("os = %q, want Windows", r.OS.Name)
	}
	if !r.IsDesktop {
		t.Error("should be desktop")
	}
}

func TestParserParseStringFirefox(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	if r.Browser.Name != "Firefox" {
		t.Errorf("browser = %q, want Firefox", r.Browser.Name)
	}
	if r.Engine.Name != "Gecko" {
		t.Errorf("engine = %q, want Gecko", r.Engine.Name)
	}
	if r.OS.Name != "Linux" {
		t.Errorf("os = %q, want Linux", r.OS.Name)
	}
}

func TestParserParseStringiPhone(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1")
	if r.OS.Name != "iOS" {
		t.Errorf("os = %q, want iOS", r.OS.Name)
	}
	if !r.IsMobile {
		t.Error("should be mobile")
	}
	if r.Device.Vendor != "Apple" {
		t.Errorf("vendor = %q, want Apple", r.Device.Vendor)
	}
}

func TestParserParseStringAndroid(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Mobile Safari/537.36")
	if r.OS.Name != "Android" {
		t.Errorf("os = %q, want Android", r.OS.Name)
	}
	if !r.IsMobile {
		t.Error("should be mobile")
	}
}

func TestParserParseStringGooglebot(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	if !r.IsBot {
		t.Error("Googlebot should be detected as bot")
	}
	if r.Bot.Name != "Googlebot" {
		t.Errorf("bot name = %q, want Googlebot", r.Bot.Name)
	}
	if r.Bot.Class != BotSearch {
		t.Errorf("bot class = %q, want search", r.Bot.Class)
	}
}

func TestParserParseStringGPTBot(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.2; +https://openai.com/gptbot)")
	if !r.IsBot {
		t.Error("GPTBot should be detected as bot")
	}
	if r.Bot.Name != "GPTBot" {
		t.Errorf("bot name = %q, want GPTBot", r.Bot.Name)
	}
	if r.Bot.Class != BotAI {
		t.Errorf("bot class = %q, want ai", r.Bot.Class)
	}
}

func TestParserParseStringCurl(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("curl/8.4.0")
	if !r.IsBot {
		t.Error("curl should be detected as bot")
	}
}

func TestParserParseStringEdge(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.2420.65")
	if r.Browser.Name != "Edge" {
		t.Errorf("browser = %q, want Edge", r.Browser.Name)
	}
	if r.Engine.Name != "Blink" {
		t.Errorf("engine = %q, want Blink", r.Engine.Name)
	}
}

func TestParserConcurrentSafety(t *testing.T) {
	p, _ := NewParser()
	uas := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 Version/17.4 Mobile/15E148 Safari/604.1",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
	}
	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func(idx int) {
			_ = p.ParseString(uas[idx%len(uas)])
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestParseOSFromComment(t *testing.T) {
	tests := []struct {
		comment string
		wantOS  string
	}{
		{"Windows NT 10.0; Win64; x64", "Windows"},
		{"Macintosh; Intel Mac OS X 10_15_7", "macOS"},
		{"Linux; Android 14; Pixel 8", "Android"},
		{"iPhone; CPU iPhone OS 17_4 like Mac OS X", "iOS"},
		{"X11; Linux x86_64", "Linux"},
		{"X11; CrOS x86_64 14541.0.0", "ChromeOS"},
	}
	for _, tt := range tests {
		os := parseOSFromComment(tt.comment)
		if os.Name != tt.wantOS {
			t.Errorf("parseOSFromComment(%q) = %q, want %q", tt.comment, os.Name, tt.wantOS)
		}
	}
}

func TestParseCPUFromComment(t *testing.T) {
	tests := []struct {
		comment  string
		wantArch string
	}{
		{"Windows NT 10.0; Win64; x64", "x86_64"},
		{"X11; Linux x86_64", "x86_64"},
		{"Linux; Android 14; Pixel 8", "arm64"},
		{"X11; Linux aarch64", "arm64"},
		{"X11; Linux i686", "x86"},
	}
	for _, tt := range tests {
		cpu := parseCPUFromComment(tt.comment)
		if cpu.Architecture != tt.wantArch {
			t.Errorf("parseCPUFromComment(%q) = %q, want %q", tt.comment, cpu.Architecture, tt.wantArch)
		}
	}
}

func TestDetectBotHeuristic(t *testing.T) {
	tests := []struct {
		ua   string
		want bool
	}{
		{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", true},
		{"crawler (+http://example.com)", true},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0 Safari/537.36", false},
		{"python-requests/2.31.0", true},
		{"curl/8.4.0", true},
	}
	for _, tt := range tests {
		bot := detectBotHeuristic(tt.ua)
		if bot.IsBot != tt.want {
			t.Errorf("detectBotHeuristic(%q) = %v, want %v", tt.ua, bot.IsBot, tt.want)
		}
	}
}

func TestGlobalParse(t *testing.T) {
	r := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0.6312.86 Safari/537.36")
	if r.Browser.Name != "Chrome" {
		t.Errorf("Parse: browser = %q, want Chrome", r.Browser.Name)
	}
}

func TestResultShortBrowser(t *testing.T) {
	r := Result{Browser: Browser{Name: "Chrome", Major: "123"}}
	if got := r.ShortBrowser(); got != "Chrome 123" {
		t.Errorf("ShortBrowser = %q, want 'Chrome 123'", got)
	}
}

func TestResultShortOS(t *testing.T) {
	r := Result{OS: OS{Name: "Windows", Major: "10"}}
	if got := r.ShortOS(); got != "Windows 10" {
		t.Errorf("ShortOS = %q, want 'Windows 10'", got)
	}
}

func TestResultDeviceClass(t *testing.T) {
	r := Result{Device: Device{Type: "mobile"}}
	if got := r.DeviceClass(); got != "mobile" {
		t.Errorf("DeviceClass = %q, want mobile", got)
	}
}

func TestResultDeviceClassUnknown(t *testing.T) {
	r := Result{}
	if got := r.DeviceClass(); got != "unknown" {
		t.Errorf("DeviceClass = %q, want unknown", got)
	}
}

type browserTestCase struct {
	UA      string `json:"ua"`
	Browser string `json:"browser"`
	Major   string `json:"major"`
	Engine  string `json:"engine"`
	OS      string `json:"os"`
	Device  string `json:"device"`
}

func TestBrowserFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/browsers.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []browserTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	p, _ := NewParser()
	for _, tc := range cases {
		t.Run(tc.Browser+"_"+tc.Device, func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.Browser.Name != tc.Browser {
				t.Errorf("browser = %q, want %q\n  UA: %s", r.Browser.Name, tc.Browser, tc.UA)
			}
			if tc.Major != "" && r.Browser.Major != tc.Major {
				t.Errorf("major = %q, want %q", r.Browser.Major, tc.Major)
			}
			if tc.Engine != "" && r.Engine.Name != tc.Engine {
				t.Errorf("engine = %q, want %q", r.Engine.Name, tc.Engine)
			}
			if tc.OS != "" && r.OS.Name != tc.OS {
				t.Errorf("os = %q, want %q", r.OS.Name, tc.OS)
			}
			if tc.Device != "" && r.DeviceClass() != tc.Device {
				t.Errorf("device = %q, want %q", r.DeviceClass(), tc.Device)
			}
		})
	}
}

type botTestCase struct {
	UA        string `json:"ua"`
	IsBot     bool   `json:"isBot"`
	BotName   string `json:"botName"`
	BotClass  string `json:"botClass"`
	BotVendor string `json:"botVendor"`
}

func TestBotFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/bots.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []botTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	p, _ := NewParser()
	for _, tc := range cases {
		name := tc.BotName
		if name == "" {
			name = "notbot"
		}
		t.Run(name, func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.IsBot != tc.IsBot {
				t.Errorf("isBot = %v, want %v\n  UA: %s", r.IsBot, tc.IsBot, tc.UA)
			}
			if tc.IsBot {
				if r.Bot.Name != tc.BotName {
					t.Errorf("botName = %q, want %q", r.Bot.Name, tc.BotName)
				}
				if string(r.Bot.Class) != tc.BotClass {
					t.Errorf("botClass = %q, want %q", r.Bot.Class, tc.BotClass)
				}
			}
		})
	}
}

func TestExpandedBrowserFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/expanded_browsers.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []browserTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	p, _ := NewParser()
	for _, tc := range cases {
		t.Run(tc.Browser+"_"+tc.Device, func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.Browser.Name != tc.Browser {
				t.Errorf("browser = %q, want %q\n  UA: %s", r.Browser.Name, tc.Browser, tc.UA)
			}
			if tc.Major != "" && r.Browser.Major != tc.Major {
				t.Errorf("major = %q, want %q", r.Browser.Major, tc.Major)
			}
			if tc.Engine != "" && r.Engine.Name != tc.Engine {
				t.Errorf("engine = %q, want %q", r.Engine.Name, tc.Engine)
			}
			if tc.OS != "" && r.OS.Name != tc.OS {
				t.Errorf("os = %q, want %q", r.OS.Name, tc.OS)
			}
			if tc.Device != "" && r.DeviceClass() != tc.Device {
				t.Errorf("device = %q, want %q", r.DeviceClass(), tc.Device)
			}
		})
	}
}

func TestExpandedBotFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/expanded_bots.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []botTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	p, _ := NewParser()
	for _, tc := range cases {
		name := tc.BotName
		if name == "" {
			name = "notbot"
		}
		t.Run(name, func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.IsBot != tc.IsBot {
				t.Errorf("isBot = %v, want %v\n  UA: %s", r.IsBot, tc.IsBot, tc.UA)
			}
			if tc.IsBot {
				if r.Bot.Name != tc.BotName {
					t.Errorf("botName = %q, want %q", r.Bot.Name, tc.BotName)
				}
				if string(r.Bot.Class) != tc.BotClass {
					t.Errorf("botClass = %q, want %q", r.Bot.Class, tc.BotClass)
				}
			}
		})
	}
}

func TestComprehensiveBrowserFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/comprehensive_browsers.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []browserTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p, _ := NewParser()
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.Browser), func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.Browser.Name != tc.Browser {
				t.Errorf("browser = %q, want %q\n  UA: %s", r.Browser.Name, tc.Browser, tc.UA)
			}
			if tc.Engine != "" && r.Engine.Name != tc.Engine {
				t.Errorf("engine = %q, want %q", r.Engine.Name, tc.Engine)
			}
			if tc.OS != "" && r.OS.Name != tc.OS {
				t.Errorf("os = %q, want %q", r.OS.Name, tc.OS)
			}
			if tc.Device != "" && r.DeviceClass() != tc.Device {
				t.Errorf("device = %q, want %q", r.DeviceClass(), tc.Device)
			}
		})
	}
}

func TestComprehensiveBotFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/comprehensive_bots.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []botTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p, _ := NewParser()
	for i, tc := range cases {
		name := tc.BotName
		if name == "" {
			name = "notbot"
		}
		t.Run(fmt.Sprintf("%d_%s", i, name), func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.IsBot != tc.IsBot {
				t.Errorf("isBot = %v, want %v\n  UA: %s", r.IsBot, tc.IsBot, tc.UA)
			}
			if tc.IsBot && tc.BotName != "" {
				if r.Bot.Name != tc.BotName {
					t.Errorf("botName = %q, want %q", r.Bot.Name, tc.BotName)
				}
				if tc.BotClass != "" && string(r.Bot.Class) != tc.BotClass {
					t.Errorf("botClass = %q, want %q", r.Bot.Class, tc.BotClass)
				}
			}
		})
	}
}

type corpusTestCase struct {
	UA       string `json:"ua"`
	Browser  string `json:"browser"`
	OS       string `json:"os"`
	IsBot    bool   `json:"isBot"`
	BotName  string `json:"botName"`
	BotClass string `json:"botClass"`
	Device   string `json:"device"`
}

func TestCorpus(t *testing.T) {
	data, err := os.ReadFile("testdata/corpus.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var cases []corpusTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p, _ := NewParser()
	passed := 0
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := p.ParseString(tc.UA)
			fail := false
			if tc.Browser != "" && r.Browser.Name != tc.Browser {
				t.Errorf("browser = %q, want %q\n  UA: %s", r.Browser.Name, tc.Browser, tc.UA)
				fail = true
			}
			if tc.OS != "" && r.OS.Name != tc.OS {
				t.Errorf("os = %q, want %q", r.OS.Name, tc.OS)
				fail = true
			}
			if r.IsBot != tc.IsBot {
				t.Errorf("isBot = %v, want %v\n  UA: %s", r.IsBot, tc.IsBot, tc.UA)
				fail = true
			}
			if tc.IsBot && tc.BotName != "" && r.Bot.Name != tc.BotName {
				t.Errorf("botName = %q, want %q", r.Bot.Name, tc.BotName)
				fail = true
			}
			if tc.IsBot && tc.BotClass != "" && string(r.Bot.Class) != tc.BotClass {
				t.Errorf("botClass = %q, want %q", r.Bot.Class, tc.BotClass)
				fail = true
			}
			if tc.Device != "" && r.DeviceClass() != tc.Device {
				t.Errorf("device = %q, want %q", r.DeviceClass(), tc.Device)
				fail = true
			}
			if !fail {
				passed++
			}
		})
	}
}

type deviceTestCase struct {
	UA           string `json:"ua"`
	DeviceType   string `json:"deviceType"`
	DeviceVendor string `json:"deviceVendor"`
}

func TestComprehensiveDeviceFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/comprehensive_devices.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var cases []deviceTestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p, _ := NewParser()
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.DeviceType), func(t *testing.T) {
			r := p.ParseString(tc.UA)
			if r.Device.Type != tc.DeviceType {
				t.Errorf("deviceType = %q, want %q\n  UA: %s", r.Device.Type, tc.DeviceType, tc.UA)
			}
			if tc.DeviceVendor != "" && r.Device.Vendor != tc.DeviceVendor {
				t.Errorf("deviceVendor = %q, want %q", r.Device.Vendor, tc.DeviceVendor)
			}
		})
	}
}
