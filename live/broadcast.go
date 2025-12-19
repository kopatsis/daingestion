package live

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type LiveEvent struct {
	Store     string
	EventCode string
	EventName string
	SessionID string
	Customer  string
	Device    string
	Location  string
	Detail    string
}

func (e LiveEvent) WithHumanizedEvent() LiveEvent {
	e.EventName = e.humanize(e.EventCode)
	return e
}

func (LiveEvent) humanize(raw string) string {
	switch raw {
	case "cart_viewed":
		return "Cart Viewed"
	case "checkout_address_info_submitted":
		return "Checkout Address Info Submitted"
	case "checkout_completed":
		return "Checkout Completed"
	case "checkout_contact_info_submitted":
		return "Checkout Contact Info Submitted"
	case "checkout_shipping_info_submitted":
		return "Checkout Shipping Info Submitted"
	case "checkout_started":
		return "Checkout Started"
	case "collection_viewed":
		return "Collection Viewed"
	case "page_viewed":
		return "Page Viewed"
	case "payment_info_submitted":
		return "Payment Info Submitted"
	case "product_added_to_cart":
		return "Product Added To Cart"
	case "product_removed_from_cart":
		return "Product Removed From Cart"
	case "product_viewed":
		return "Product Viewed"
	case "search_submitted":
		return "Search Submitted"
	default:
		return raw
	}
}

func PublishEvent(ctx context.Context, rdb *redis.Client, e LiveEvent) error {
	ch := fmt.Sprintf("events:%s:%s", e.Store, e.EventCode)
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return rdb.Publish(context.TODO(), ch, b).Err()
}
