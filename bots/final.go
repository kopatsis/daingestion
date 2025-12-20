package bots

import "dmd/models"

type BotLevel int

const (
	NotBot BotLevel = iota
	MaybeBot
	LikelyBot
)

func EvaluateBot(
	req models.RequestSignals,
	env models.BotSignals,
	asnIsDatacenter bool,
	uaDeviceDetectorBot bool,
) BotLevel {
	score := 0

	if uaDeviceDetectorBot {
		return LikelyBot
	}
	if req.IsBotCrawlerDetect {
		return LikelyBot
	}

	if asnIsDatacenter {
		score += 3
	}

	if req.MethodInvalid {
		score += 3
	}
	if req.InvalidHTTPVersion {
		score += 2
	}
	if req.ConnectionClose {
		score += 1
	}
	if req.XFFMalformed {
		score += 3
	} else if req.XFFPrivate {
		score += 2
	}

	if env.ViewportImpossible {
		score += 2
	}
	if env.ViewportContradiction {
		score += 1
	}
	if env.NavigatorCookieFalse {
		score += 1
	}
	if env.NavigatorLangEmpty {
		score += 1
	}
	if env.NavigatorLangsEmpty {
		score += 1
	}
	if env.NavigatorUAEmpty {
		score += 2
	}

	if score >= 6 {
		return LikelyBot
	}
	if score >= 2 {
		return MaybeBot
	}
	return NotBot
}
