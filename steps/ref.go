package steps

import (
	"net/url"
	"strings"
)

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

func ParseReferrer(raw string) Referrer {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return Referrer{Raw: raw}
	}

	host := u.Hostname()
	port := u.Port()
	origin := u.Scheme + "://" + u.Host
	q := strings.TrimPrefix(u.RawQuery, "")
	f := strings.TrimPrefix(u.Fragment, "")

	search, engine := SearchEngine(host)

	return Referrer{
		Exists:         false,
		Raw:            raw,
		Protocol:       u.Scheme,
		Hostname:       host,
		Port:           port,
		Path:           u.EscapedPath(),
		Query:          q,
		Fragment:       f,
		Origin:         origin,
		DomainOnly:     host,
		IsSearchEngine: search,
		SearchEngine:   engine,
	}
}

func SearchEngine(host string) (bool, string) {
	h := strings.ToLower(host)

	switch {
	case strings.Contains(h, "google."):
		return true, "google"
	case strings.Contains(h, "bing."):
		return true, "bing"
	case strings.Contains(h, "yahoo."):
		return true, "yahoo"
	case strings.Contains(h, "duckduckgo.") || strings.Contains(h, "ddg.gg"):
		return true, "duckduckgo"
	case strings.Contains(h, "baidu."):
		return true, "baidu"
	case strings.Contains(h, "yandex."):
		return true, "yandex"
	case strings.Contains(h, "naver."):
		return true, "naver"
	case strings.Contains(h, "seznam."):
		return true, "seznam"
	case strings.Contains(h, "ecosia."):
		return true, "ecosia"
	case strings.Contains(h, "startpage.") || strings.Contains(h, "ixquick."):
		return true, "startpage"
	case strings.Contains(h, "qwant."):
		return true, "qwant"
	}

	return false, ""
}
