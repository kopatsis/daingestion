package models

import "encoding/json"

type IngestEvent struct {
	Event struct {
		Name      string `json:"name"`
		Timestamp int64  `json:"timestamp"`
		ClientID  string `json:"clientID"`
		Context   struct {
			Window struct {
				InnerHeight int `json:"innerHeight"`
				InnerWidth  int `json:"innerWidth"`
				Screen      struct {
					Height int `json:"height"`
					Width  int `json:"width"`
				} `json:"screen"`
				Rest json.RawMessage `json:"-"`
			} `json:"window"`
			Navigator struct {
				CookieEnabled bool            `json:"cookieEnabled"`
				UserAgent     string          `json:"userAgent"`
				Rest          json.RawMessage `json:"-"`
			} `json:"navigator"`
			Document struct {
				Referrer string `json:"referrer"`
				Location struct {
					Href string `json:"href"`
				} `json:"location"`
				Rest json.RawMessage `json:"-"`
			} `json:"document"`
			Rest json.RawMessage `json:"-"`
		} `json:"context"`
		Data json.RawMessage `json:"data"`
		Rest json.RawMessage `json:"-"`
	} `json:"event"`

	Init struct {
		Data struct {
			Shop struct {
				Name            string          `json:"name"`
				MyShopifyDomain string          `json:"myshopifyDomain"`
				Rest            json.RawMessage `json:"-"`
			} `json:"shop"`
			Rest json.RawMessage `json:"-"`
		} `json:"data"`
		Rest json.RawMessage `json:"-"`
	} `json:"init"`

	Session struct {
		Current  string `json:"current"`
		Previous string `json:"previous"`
	} `json:"session"`

	Time int64 `json:"time"`
}
