package models

type UAInfo struct {
	Model           string
	Brand           string
	Type            string
	OSVersion       string
	OSShortName     string
	OSName          string
	OSPlatform      string
	ClientType      string
	ClientName      string
	ClientVersion   string
	ClientShortName string
	ClientEngine    string
	ClientEngineVer string
	BotName         string
	BotCategory     string
	BotProducerName string
	BotProducerURL  string
	BotURL          string
	IsBot           bool
	IsMobile        bool
	IsDesktop       bool
	IsTouch         bool
}

type GeoData struct {
	IP              string
	CountryISO      string
	CountryName     string
	SubdivisionISO  string
	SubdivisionName string
	CityName        string
	PostalCode      string
	Latitude        float64
	Longitude       float64
	AccuracyRadius  uint16
	ASN             uint
	ASNOrg          string
	DataCenter      string
}

type UTM struct {
	Source   string
	Medium   string
	Campaign string
	Term     string
	Content  string
	All      map[string]string
}

type Referrer struct {
	Exists         bool
	Raw            string
	Protocol       string
	Hostname       string
	Port           string
	Path           string
	Query          string
	Fragment       string
	Origin         string
	DomainOnly     string
	IsSearchEngine bool
	SearchEngine   string
}

type ScreenBuckets struct {
	InnerWidthBucket   string
	InnerHeightBucket  string
	ScreenWidthBucket  string
	ScreenHeightBucket string
}

type RequestSignals struct {
	ReferrerEmpty      bool
	ReferrerMalformed  bool
	IsBotCrawlerDetect bool
	ConnectionClose    bool
	MethodInvalid      bool
	InvalidHTTPVersion bool
	XFFEmpty           bool
	XFFPrivate         bool
	XFFMalformed       bool
}

type BotSignals struct {
	ViewportImpossible    bool
	ViewportContradiction bool
	NavigatorCookieFalse  bool
	NavigatorLangEmpty    bool
	NavigatorLangsEmpty   bool
	NavigatorUAEmpty      bool
}
