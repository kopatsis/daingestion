package live

import (
	"context"
	"dmd/steps"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func AddActive(rdb *redis.Client, tenant, dim, id, sessionID string) {
	now := time.Now()
	b := now.Unix() / 15
	k := fmt.Sprintf("active:%s:%s:%s:%d", tenant, dim, id, b)
	rdb.SAdd(context.TODO(), k, sessionID)
	rdb.Expire(context.TODO(), k, 4*time.Minute+15*time.Second)
}

func SetSessionMetadata(ctx context.Context, rdb *redis.Client, sessionID string, g steps.GeoData) error {
	key := "session:" + sessionID
	_, err := rdb.HSet(ctx, key, map[string]interface{}{
		"IP":              g.IP,
		"CountryISO":      g.CountryISO,
		"CountryName":     g.CountryName,
		"SubdivisionISO":  g.SubdivisionISO,
		"SubdivisionName": g.SubdivisionName,
		"CityName":        g.CityName,
		"PostalCode":      g.PostalCode,
		"Latitude":        strconv.FormatFloat(g.Latitude, 'f', -1, 64),
		"Longitude":       strconv.FormatFloat(g.Longitude, 'f', -1, 64),
		"AccuracyRadius":  strconv.Itoa(int(g.AccuracyRadius)),
		"ASN":             strconv.Itoa(int(g.ASN)),
		"ASNOrg":          g.ASNOrg,
	}).Result()
	if err != nil {
		return err
	}
	return rdb.Expire(ctx, key, 300000000000).Err()
}
