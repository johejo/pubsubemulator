package pubsubemulator_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/johejo/pubsubemulator"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	e, err := pubsubemulator.New(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e.Host())
	t.Setenv("PUBSUB_EMULATOR_HOST", e.Host())
	client, err := pubsub.NewClient(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	topic, err := client.CreateTopic(ctx, "test-topic")
	if err != nil {
		t.Fatal(err)
	}

	sub, err := client.CreateSubscription(ctx, "test-sub", pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		t.Fatal(err)
	}

	finish := make(chan struct{})
	go func() {
		if err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			t.Log("receive", string(m.Data))
			close(finish)

		}); err != nil {
			panic(err)
		}
	}()

	result := topic.Publish(ctx, &pubsub.Message{Data: []byte("test")})
	serverID, err := result.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("serverID", serverID)

	select {
	case <-finish:
	case <-ctx.Done():
		t.Error(ctx.Err())
	}

	if err := e.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}
