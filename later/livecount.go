package later

import (
	"context"
	"dmd/models"
	"fmt"
	"strconv"
	"time"

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

func GetActiveLast16Buckets(rdb *redis.Client, store string) ([]string, error) {
	now := time.Now().Unix()
	currentBucket := now / 15

	if now%15 != 0 {
		currentBucket--
	}

	startBucket := currentBucket - 15

	ctx := context.TODO()
	combined := make(map[string]struct{})

	for b := startBucket; b <= currentBucket; b++ {
		key := fmt.Sprintf("active:%s:%d", store, b)
		members, err := rdb.SMembers(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		for _, m := range members {
			combined[m] = struct{}{}
		}
	}

	out := make([]string, 0, len(combined))
	for k := range combined {
		out = append(out, k)
	}

	return out, nil
}

func GetStates(
	ctx context.Context,
	rdb *redis.Client,
	sessionIDs []string,
) (map[string]models.SessionActiveState, error) {

	pipe := rdb.Pipeline()
	cmds := make(map[string]*redis.MapStringStringCmd, len(sessionIDs))

	for _, sid := range sessionIDs {
		key := "session:" + sid
		cmds[sid] = pipe.HGetAll(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	out := make(map[string]models.SessionActiveState, len(sessionIDs))

	for sid, cmd := range cmds {
		m := cmd.Val()
		if len(m) == 0 {
			continue
		}

		asn, _ := strconv.ParseUint(m["ASN"], 10, 64)
		isBot, _ := strconv.ParseBool(m["IsBot"])
		isLoggedIn, _ := strconv.ParseBool(m["IsLoggedIn"])
		hasPrev, _ := strconv.ParseBool(m["HasPreviousPurchase"])
		isViewingProduct, _ := strconv.ParseBool(m["IsViewingProduct"])
		isViewingCollection, _ := strconv.ParseBool(m["IsViewingCollection"])
		hasActiveCart, _ := strconv.ParseBool(m["HasActiveCart"])
		isViewingCart, _ := strconv.ParseBool(m["IsViewingCart"])
		isInCheckout, _ := strconv.ParseBool(m["IsInCheckout"])
		lat, _ := strconv.ParseFloat(m["Latitude"], 64)
		lon, _ := strconv.ParseFloat(m["Longitude"], 64)

		out[sid] = models.SessionActiveState{
			Country:   m["Country"],
			Region:    m["Region"],
			City:      m["City"],
			Latitude:  lat,
			Longitude: lon,
			Timezone:  m["Timezone"],

			IsBot:       isBot,
			ASN:         uint(asn),
			ASNProvider: m["ASNProvider"],

			IsLoggedIn:          isLoggedIn,
			HasPreviousPurchase: hasPrev,

			DeviceType:  m["DeviceType"],
			DeviceBrand: m["DeviceBrand"],
			OSName:      m["OSName"],
			BrowserName: m["BrowserName"],
			Language:    m["Language"],

			RefDomain: m["RefDomain"],
			UTMSource: m["UTMSource"],

			RouteType: m["RouteType"],
			Route:     m["Route"],
			FullURL:   m["FullURL"],

			IsViewingProduct: isViewingProduct,
			ActiveProductID:  m["ActiveProductID"],
			ActiveVariantID:  m["ActiveVariantID"],

			IsViewingCollection: isViewingCollection,
			ActiveCollectionID:  m["ActiveCollectionID"],

			HasActiveCart: hasActiveCart,
			IsViewingCart: isViewingCart,
			IsInCheckout:  isInCheckout,
		}
	}

	return out, nil
}

func GetSpecificSessionFields(
	ctx context.Context,
	rdb *redis.Client,
	sessionIDs []string,
	fields []string,
) (map[string]map[string]string, error) {

	pipe := rdb.Pipeline()
	cmds := make(map[string]*redis.SliceCmd, len(sessionIDs))

	for _, sid := range sessionIDs {
		key := "session:" + sid
		cmds[sid] = pipe.HMGet(ctx, key, fields...)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	out := make(map[string]map[string]string, len(sessionIDs))

	for sid, cmd := range cmds {
		vals := cmd.Val()
		if len(vals) == 0 {
			continue
		}

		m := make(map[string]string, len(fields))
		for i, v := range vals {
			if v == nil {
				continue
			}
			m[fields[i]] = v.(string)
		}

		if len(m) > 0 {
			out[sid] = m
		}
	}

	return out, nil
}
