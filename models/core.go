package models

import "encoding/json"

type IngestEvent struct {
	Event struct {
		Name      string `json:"name" validate:"required"`
		Timestamp int64  `json:"timestamp" validate:"required"`
		ClientID  string `json:"clientID" validate:"required"`
		Context   struct {
			Window struct {
				InnerHeight int `json:"innerHeight" validate:"required"`
				InnerWidth  int `json:"innerWidth" validate:"required"`
				Screen      struct {
					Height int `json:"height" validate:"required"`
					Width  int `json:"width" validate:"required"`
				} `json:"screen" validate:"required"`
				Rest json.RawMessage `json:"-"`
			} `json:"window" validate:"required"`
			Navigator Navigator `json:"navigator" validate:"required"`
			Document  struct {
				Referrer string `json:"referrer" validate:"required"`
				Location struct {
					Href   string `json:"href" validate:"required"`
					Search string `json:"search" validate:"required"`
				} `json:"location" validate:"required"`
				Rest json.RawMessage `json:"-"`
			} `json:"document" validate:"required"`
			Rest json.RawMessage `json:"-"`
		} `json:"context" validate:"required"`
		Data json.RawMessage `json:"data" validate:"required"`
		Rest json.RawMessage `json:"-"`
	} `json:"event" validate:"required"`

	Init struct {
		Data struct {
			Shop struct {
				Name            string          `json:"name" validate:"required"`
				MyShopifyDomain string          `json:"myshopifyDomain" validate:"required"`
				Rest            json.RawMessage `json:"-"`
			} `json:"shop" validate:"required"`
			Rest json.RawMessage `json:"-"`
		} `json:"data" validate:"required"`
		Rest json.RawMessage `json:"-"`
	} `json:"init" validate:"required"`

	Time int64 `json:"time" validate:"required"`
}

type Navigator struct {
	CookieEnabled bool     `json:"cookieEnabled" validate:"required"`
	Language      string   `json:"userAgent" validate:"required"`
	Languages     []string `json:"languages" validate:"required"`
	UserAgent     string   `json:"language" validate:"required"`
}
