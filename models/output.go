package models

type Output struct {
	EventName string
	EventID   string

	EventTime int64
	TimeIn    int64
	TimeOut   int64

	ShopName   string
	ShopDomain string

	ClientID   string
	SessionID  string
	CustomerID string
	LoggedIn   bool
	IP         string
	IPHash     string

	SessionStatus    string
	PageType         string
	PreviousPurchase bool
	BotScore         int

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
