package output

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub/v2"
)

type Output struct {
	EventName string `json:"event_name"`
	Timestamp int64  `json:"timestamp"`
	Data      any    `json:"data"`
}

func NewPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID)
}

func PublishOutput(ctx context.Context, client *pubsub.Client, topicID string, out Output) error {
	b, err := json.Marshal(out)
	if err != nil {
		return err
	}
	p := client.Publisher(topicID)
	r := p.Publish(ctx, &pubsub.Message{Data: b})
	_, err = r.Get(ctx)
	return err
}
