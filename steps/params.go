package steps

import (
	"dmd/models"
	"net/url"
	"strings"
)

func ParseUTM(search string) models.UTM {
	u := models.UTM{All: make(map[string]string)}
	s := strings.TrimPrefix(search, "?")
	if s == "" {
		return u
	}

	m, err := url.ParseQuery(s)
	if err != nil {
		return u
	}

	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		lk := strings.ToLower(k)

		switch lk {
		case "utm_source":
			u.Source = val
		case "utm_medium":
			u.Medium = val
		case "utm_campaign":
			u.Campaign = val
		case "utm_term":
			u.Term = val
		case "utm_content":
			u.Content = val
		default:
			if strings.HasPrefix(lk, "utm_") {
				u.All[lk] = val
			}
		}
	}

	return u
}

func ParseNonUTMParams(search string) map[string]string {
	out := make(map[string]string)

	s := strings.TrimPrefix(search, "?")
	if s == "" {
		return out
	}

	m, err := url.ParseQuery(s)
	if err != nil {
		return out
	}

	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		lk := strings.ToLower(k)

		if strings.HasPrefix(lk, "utm") {
			continue
		}

		out[k] = v[0]
	}

	return out
}
