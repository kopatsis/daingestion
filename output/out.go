package output

import (
	"context"
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
		return err
	}
	p := client.Publisher(topicID)
	r := p.Publish(ctx, &pubsub.Message{Data: b})
	_, err = r.Get(ctx)
	return err
}
