package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/gamebtc/devicedetector"
	"github.com/oschwald/maxminddb-golang/v2"
)

type DataCenterASNs struct {
	set map[uint]struct{}
}

var (
	asnList DataCenterASNs
	geoDB   *maxminddb.Reader
)

func LoadASNs(path string) (DataCenterASNs, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return DataCenterASNs{}, err
	}
	s := DataCenterASNs{set: map[uint]struct{}{}}
	start := 0
	for i := 0; i <= len(b); i++ {
		if i == len(b) || b[i] == '\n' {
			if i > start {
				var n uint
				parseNum(b[start:i], &n)
				s.set[n] = struct{}{}
			}
			start = i + 1
		}
	}
	return s, nil
}

func parseNum(b []byte, out *uint) {
	var n uint
	for _, c := range b {
		if c >= '0' && c <= '9' {
			n = n*10 + uint(c-'0')
		}
	}
	*out = n
}

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

func main() {
	a, err := LoadASNs("files/asn.txt")
	if err != nil {
		panic(err)
	}
	asnList = a

	db, err := maxminddb.Open("files/GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	geoDB = db

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.ListenAndServe(":8080", nil)
}
