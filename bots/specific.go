package bots

import (
	"dmd/models"
	"net/http"
)

func EvaluateSpecific(r *http.Request, ref string, nav models.Navigator, innerW, innerH, screenW, screenH int, shopDomain string) models.BotSignals {
	s := models.BotSignals{}

	if innerW <= 0 || innerH <= 0 || screenW <= 0 || screenH <= 0 {
		s.ViewportImpossible = true
	}
	if innerW > screenW || innerH > screenH {
		s.ViewportContradiction = true
	}

	if !nav.CookieEnabled {
		s.NavigatorCookieFalse = true
	}
	if nav.Language == "" {
		s.NavigatorLangEmpty = true
	}
	if len(nav.Languages) == 0 {
		s.NavigatorLangsEmpty = true
	}
	if nav.UserAgent == "" {
		s.NavigatorUAEmpty = true
	}

	return s
}
