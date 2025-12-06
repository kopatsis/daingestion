package steps

import (
	"fmt"
	"log"

	. "github.com/gamebtc/devicedetector"
)

func L() {
	dd, err := NewDeviceDetector("regexes")
	if err != nil {
		log.Fatal(err)
	}
	userAgent := `Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1`
	info := dd.Parse(userAgent)

	fmt.Println(info.Model) // iPhone
	fmt.Println(info.Brand) // AP
	fmt.Println(info.Type)  // smartphone

	os := info.GetOs()        //
	fmt.Println(os.Version)   // 11.0
	fmt.Println(os.ShortName) // IOS
	fmt.Println(os.Name)      // iOS
	fmt.Println(os.Platform)  //

	client := info.GetClient()
	fmt.Println(client.Type)    // browser
	fmt.Println(client.Name)    // Mobile Safari
	fmt.Println(client.Version) // 11.0

	if client.Type == `browser` {
		fmt.Println(client.ShortName)     // MF
		fmt.Println(client.Engine)        // WebKit
		fmt.Println(client.EngineVersion) // 604.1.38
	}

	bot := info.GetBot()
	if bot != nil {
		fmt.Println(bot.Name)
		//.................
	}
}
