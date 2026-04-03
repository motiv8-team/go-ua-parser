package uax

// Browser describes the detected web browser.
type Browser struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Major   string `json:"major,omitempty"`
	Family  string `json:"family,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// Engine describes the browser rendering engine.
type Engine struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// OS describes the operating system.
type OS struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Major   string `json:"major,omitempty"`
	Minor   string `json:"minor,omitempty"`
	Patch   string `json:"patch,omitempty"`
}

// CPU describes the processor architecture.
type CPU struct {
	Architecture string `json:"architecture,omitempty"`
	Bits         int    `json:"bits,omitempty"`
}

// Device describes the physical device.
type Device struct {
	Type      string `json:"type,omitempty"`
	Vendor    string `json:"vendor,omitempty"`
	Model     string `json:"model,omitempty"`
	IsTouch   bool   `json:"isTouch,omitempty"`
	IsPhone   bool   `json:"isPhone,omitempty"`
	IsTablet  bool   `json:"isTablet,omitempty"`
	IsDesktop bool   `json:"isDesktop,omitempty"`
	IsTV      bool   `json:"isTV,omitempty"`
}

// App describes an in-app browser or wrapper application.
type App struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Kind        string `json:"kind,omitempty"`
}

// Result holds the complete parsed output from a User-Agent string.
type Result struct {
	UAString  string   `json:"uaString,omitempty"`
	Browser   Browser  `json:"browser"`
	Engine    Engine   `json:"engine"`
	OS        OS       `json:"os"`
	CPU       CPU      `json:"cpu"`
	Device    Device   `json:"device"`
	App       App      `json:"app"`
	Bot       Bot      `json:"bot"`
	IsMobile  bool     `json:"isMobile,omitempty"`
	IsDesktop bool     `json:"isDesktop,omitempty"`
	IsTablet  bool     `json:"isTablet,omitempty"`
	IsBot     bool     `json:"isBot,omitempty"`
	IsCrawler bool     `json:"isCrawler,omitempty"`
	IsInApp   bool     `json:"isInApp,omitempty"`
}
