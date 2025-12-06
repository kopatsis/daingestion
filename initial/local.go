package initial

import (
	"os"

	"github.com/oschwald/maxminddb-golang/v2"
)

type DataCenterASNs struct {
	set map[uint]struct{}
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

func LoadASNs() (DataCenterASNs, error) {
	b, err := os.ReadFile("files/asn.txt")
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
