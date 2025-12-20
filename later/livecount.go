package later

import (
	"context"
	"dmd/models"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func GetSessionField(ctx context.Context, rdb *redis.Client, sessionID, field string) (string, error) {
	key := "session:" + sessionID
	v, err := rdb.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return v, err
}

func GetSessionAll(ctx context.Context, rdb *redis.Client, sessionID string) (models.GeoData, error) {
	key := "session:" + sessionID
	m, err := rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return models.GeoData{}, nil
	}
	if err != nil {
		return models.GeoData{}, err
	}

	lat, _ := strconv.ParseFloat(m["Latitude"], 64)
	lon, _ := strconv.ParseFloat(m["Longitude"], 64)
	acc, _ := strconv.ParseUint(m["AccuracyRadius"], 10, 16)
	asn, _ := strconv.ParseUint(m["ASN"], 10, 32)

	return models.GeoData{
		IP:              m["IP"],
		CountryISO:      m["CountryISO"],
		CountryName:     m["CountryName"],
		SubdivisionISO:  m["SubdivisionISO"],
		SubdivisionName: m["SubdivisionName"],
		CityName:        m["CityName"],
		PostalCode:      m["PostalCode"],
		Latitude:        lat,
		Longitude:       lon,
		AccuracyRadius:  uint16(acc),
		ASN:             uint(asn),
		ASNOrg:          m["ASNOrg"],
	}, nil
}
