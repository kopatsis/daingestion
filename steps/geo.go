package steps

import (
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

type cityRecord struct {
	Country struct {
		ISOCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`

	Subdivisions []struct {
		ISOCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`

	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`

	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`

	Location struct {
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
}

type asnRecord struct {
	ASN uint   `maxminddb:"autonomous_system_number"`
	Org string `maxminddb:"autonomous_system_organization"`
}

type GeoData struct {
	IP              string
	CountryISO      string
	CountryName     string
	SubdivisionISO  string
	SubdivisionName string
	CityName        string
	PostalCode      string
	Latitude        float64
	Longitude       float64
	AccuracyRadius  uint16
	ASN             uint
	ASNOrg          string
}

func ExtractGeo(ipStr string, cityDB *maxminddb.Reader, asnDB *maxminddb.Reader) GeoData {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return GeoData{IP: ipStr}
	}

	var c cityRecord
	var a asnRecord

	_ = cityDB.Lookup(ip).Decode(&c)
	_ = asnDB.Lookup(ip).Decode(&a)

	subISO := ""
	subName := ""
	if len(c.Subdivisions) > 0 {
		subISO = c.Subdivisions[0].ISOCode
		subName = c.Subdivisions[0].Names["en"]
	}

	return GeoData{
		IP:              ipStr,
		CountryISO:      c.Country.ISOCode,
		CountryName:     c.Country.Names["en"],
		SubdivisionISO:  subISO,
		SubdivisionName: subName,
		CityName:        c.City.Names["en"],
		PostalCode:      c.Postal.Code,
		Latitude:        c.Location.Latitude,
		Longitude:       c.Location.Longitude,
		AccuracyRadius:  c.Location.AccuracyRadius,
		ASN:             a.ASN,
		ASNOrg:          a.Org,
	}
}
