package uax

// BotClass categorizes the type of automated agent.
type BotClass string

const (
	BotNone    BotClass = ""         // BotNone indicates no bot was detected.
	BotSearch  BotClass = "search"   // BotSearch identifies a search engine crawler.
	BotSocial  BotClass = "social"   // BotSocial identifies a social media crawler.
	BotMonitor BotClass = "monitor"  // BotMonitor identifies an uptime or performance monitor.
	BotScraper BotClass = "scraper"  // BotScraper identifies a generic web scraper.
	BotAI      BotClass = "ai"       // BotAI identifies an AI/LLM training or inference crawler.
	BotSEO     BotClass = "seo-tool" // BotSEO identifies an SEO analysis tool.
	BotOther   BotClass = "other"    // BotOther covers automated agents that don't fit other categories.
)

// Bot describes a detected automated agent (crawler, bot, AI agent, etc.).
type Bot struct {
	IsBot      bool     `json:"isBot,omitempty"`
	Class      BotClass `json:"class,omitempty"`
	Name       string   `json:"name,omitempty"`
	Version    string   `json:"version,omitempty"`
	Vendor     string   `json:"vendor,omitempty"`
	IsVerified bool    `json:"isVerified,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}
