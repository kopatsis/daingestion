package steps

import (
	"github.com/gamebtc/devicedetector"
)

type UAInfo struct {
	Model           string
	Brand           string
	Type            string
	OSVersion       string
	OSShortName     string
	OSName          string
	OSPlatform      string
	ClientType      string
	ClientName      string
	ClientVersion   string
	ClientShortName string
	ClientEngine    string
	ClientEngineVer string
	BotName         string
	BotCategory     string
	BotProducerName string
	BotProducerURL  string
	BotURL          string
	IsBot           bool
	IsMobile        bool
	IsDesktop       bool
	IsTouch         bool
}

func ParseUA(s string) UAInfo {
	dd, _ := devicedetector.NewDeviceDetector("regexes")
	info := dd.Parse(s)

	os := info.GetOs()
	client := info.GetClient()
	bot := info.GetBot()

	out := UAInfo{
		Model:         info.Model,
		Brand:         info.Brand,
		Type:          info.Type,
		OSVersion:     os.Version,
		OSShortName:   os.ShortName,
		OSName:        os.Name,
		OSPlatform:    os.Platform,
		ClientType:    client.Type,
		ClientName:    client.Name,
		ClientVersion: client.Version,
		IsBot:         info.IsBot(),
		IsMobile:      info.IsMobile(),
		IsDesktop:     info.IsDesktop(),
		IsTouch:       info.IsTouchEnabled(),
	}

	if client.Type == "browser" {
		out.ClientShortName = client.ShortName
		out.ClientEngine = client.Engine
		out.ClientEngineVer = client.EngineVersion
	}

	if bot != nil {
		out.BotName = bot.Name
		out.BotCategory = bot.Category
		out.BotProducerName = bot.Producer.Name
		out.BotProducerURL = bot.Producer.Url
		out.BotURL = bot.Url
	}

	return out
}
