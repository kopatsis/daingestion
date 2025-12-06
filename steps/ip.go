package steps

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) (string, string) {
	ip := ""
	h := r.Header.Get("X-Forwarded-For")
	if h != "" {
		parts := strings.Split(h, ",")
		c := strings.TrimSpace(parts[0])
		if net.ParseIP(c) != nil {
			ip = c
		}
	}
	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && net.ParseIP(host) != nil {
			ip = host
		}
	}
	if ip == "" {
		return "", ""
	}
	hx := sha256.Sum256([]byte(ip))
	return ip, hex.EncodeToString(hx[:])
}
