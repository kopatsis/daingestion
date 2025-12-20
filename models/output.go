package models

type IngestedEvent struct {
	EventName      string
	EventTimestamp int64
	ShopID         string
	ShopDomain     string

	SessionID        string
	SessionStatus    string
	PageType         string
	LoggedIn         bool
	PreviousPurchase bool
	BotScore         int
	Datacenter       string
	IP               string
	IPHash           string

	Params map[string]string

	UA             UAInfo
	Geo            GeoData
	UTM            UTM
	Referrer       Referrer
	Screen         ScreenBuckets
	RequestSignals RequestSignals
	BotSignals     BotSignals

	RawShopify IngestEvent
}
