package steps

import (
	"dmd/models"

	"github.com/gamebtc/devicedetector"
)

func ParseUA(s string) models.UAInfo {
	dd, _ := devicedetector.NewDeviceDetector("regexes")
	info := dd.Parse(s)

	os := info.GetOs()
	client := info.GetClient()
	bot := info.GetBot()

	out := models.UAInfo{
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
