package bots

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/x-way/crawlerdetect"
)

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

func ExtractSignals(r *http.Request, rawRef, rawUA string) RequestSignals {
	connection := r.Header.Get("Connection")
	xff := r.Header.Get("X-Forwarded-For")
	httpVer := r.Proto
	method := r.Method

	cd := crawlerdetect.New()

	s := RequestSignals{}

	if rawRef == "" {
		s.ReferrerEmpty = true
	} else {
		u, err := url.Parse(rawRef)
		if err != nil || u.Scheme == "" || u.Host == "" {
			s.ReferrerMalformed = true
		}
	}

	if cd.IsCrawler(rawUA) {
		s.IsBotCrawlerDetect = true
	}

	if strings.ToLower(connection) == "close" {
		s.ConnectionClose = true
	}

	if xff == "" {
		s.XFFEmpty = true
	} else {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		parsed := net.ParseIP(ip)
		if parsed == nil {
			s.XFFMalformed = true
		} else if parsed.IsPrivate() {
			s.XFFPrivate = true
		}
	}

	if strings.Contains(httpVer, "HTTP/1.0") || strings.Contains(httpVer, "HTTP/0") {
		s.InvalidHTTPVersion = true
	}

	if method == "HEAD" || method == "OPTIONS" {
		s.MethodInvalid = true
	}

	return s
}
