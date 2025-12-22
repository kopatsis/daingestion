package live

import (
	"context"
	"dmd/logging"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type LiveEvent struct {
	Store      string
	EventCode  string
	EventName  string
	SessionID  string
	Customer   string
	Device     string
	Location   string
	URL        string
	Search     *string
	Product    *string
	Collection *string
	OrderID    *string
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

func PublishEvent(ctx context.Context, rdb *redis.Client, e LiveEvent, reqID string) error {
	ch := fmt.Sprintf("events:%s:%s", e.Store, e.EventCode)
	b, err := json.Marshal(e)
	if err != nil {
		logging.LogError(
			"ERROR",
			"live_broadcast_failed",
			"redis",
			e.Store,
			e.EventName,
			reqID,
			false,
			"failed to marshal broadcast live event",
		)

		return err
	}
	if err := rdb.Publish(context.TODO(), ch, b).Err(); err != nil {
		logging.LogError(
			"ERROR",
			"live_broadcast_failed",
			"redis",
			e.Store,
			e.EventName,
			reqID,
			false,
			"failed to broadcast live event",
		)
		return err
	}
	return nil
}

func StreamEventsSSE(
	rdb *redis.Client,
	store string,
	eventCode string,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "stream unsupported", http.StatusInternalServerError)
			return
		}

		ch := fmt.Sprintf("events:%s:%s", store, eventCode)
		sub := rdb.Subscribe(ctx, ch)
		_, err := sub.Receive(ctx)
		if err != nil {
			return
		}
		defer sub.Close()

		msgCh := sub.Channel()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgCh:
				if msg == nil {
					return
				}

				var e LiveEvent
				if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
					continue
				}

				b, err := json.Marshal(e)
				if err != nil {
					continue
				}

				fmt.Fprintf(w, "data: %s\n\n", b)
				flusher.Flush()
			}
		}
	}
}
