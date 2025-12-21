package output

import (
	"context"
	"dmd/logging"
	"dmd/models"
	"encoding/json"

	"cloud.google.com/go/pubsub/v2"
)

func NewPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID)
}

func PublishOutput(ctx context.Context, client *pubsub.Client, topicID string, out models.Output) error {
	b, err := json.Marshal(out)
	if err != nil {
		logging.LogError(
			"CRITICAL",
			"pubsub_marshal_failed",
			"pubsub",
			out.ShopDomain,
			out.EventName,
			out.RequestID,
			true,
			"failed to marshal event to pubsub",
		)
		return err
	}
	p := client.Publisher(topicID)
	r := p.Publish(ctx, &pubsub.Message{Data: b})
	if _, err := r.Get(ctx); err != nil {
		logging.LogError(
			"CRITICAL",
			"pubsub_publish_failed",
			"pubsub",
			out.ShopDomain,
			out.EventName,
			out.RequestID,
			true,
			"failed to publish event to pubsub",
		)
		return err
	}
	return nil
}
