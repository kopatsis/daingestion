package initial

import (
	"bytes"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang/v2"
)

type DataCenterASNs struct {
	Orgs map[uint]string
}

func LoadASNs() (DataCenterASNs, error) {
	b, err := os.ReadFile("files/asn.txt")
	if err != nil {
		return DataCenterASNs{}, err
	}

	out := DataCenterASNs{Orgs: make(map[uint]string)}

	start := 0
	for i := 0; i <= len(b); i++ {
		if i != len(b) && b[i] != '\n' {
			continue
		}

		if i > start {
			line := b[start:i]
			space := bytes.IndexByte(line, ' ')
			if space > 0 {
				asPart := line[:space]
				namePart := strings.TrimSpace(string(line[space+1:]))

				var n uint
				for _, c := range asPart {
					if c >= '0' && c <= '9' {
						n = n*10 + uint(c-'0')
					}
				}

				out.Orgs[n] = namePart
			}
		}

		start = i + 1
	}

	return out, nil
}

func LoadGeoLite() (*maxminddb.Reader, *maxminddb.Reader, error) {
	cityDB, err := maxminddb.Open("files/GeoLite2-City.mmdb")
	if err != nil {
		return nil, nil, err
	}

	asnDB, err := maxminddb.Open("files/GeoLite2-ASN.mmdb")
	if err != nil {
		cityDB.Close()
		return nil, nil, err
	}

	return cityDB, asnDB, nil
}
