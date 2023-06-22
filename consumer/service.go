package consumer

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/segmentio/kafka-go"
)

type config struct {
	Topic   string
	Brokers []string
	GroupID string
}

type Service interface{}

type consumer struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]
	reader *kafka.Reader
}

func (c *consumer) Init(ctx context.Context) error {
	c.Logger().
		With("brokers", c.Config().Brokers).
		With("topic", c.Config().Topic).
		With("group-id", c.Config().GroupID).
		Info("initializing consumer ...")

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Topic:   c.Config().Topic,
		Brokers: c.Config().Brokers,
		GroupID: c.Config().GroupID,
	})

	go c.listen(ctx)

	return nil
}

func (c *consumer) listen(ctx context.Context) {
	logger := c.Logger()
listener:
	for {
		select {
		case <-ctx.Done():
			break listener
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				logger.ErrorCtx(ctx, "failed to read message from kafka", "reason", err)
				continue
			}
			logger.Info("received message", "offset", msg.Offset, "message", string(msg.Value))
		}
	}

	if err := c.reader.Close(); err != nil {
		logger.ErrorCtx(ctx, "failed to close kafka reader")
	}
}
