package consumer

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/segmentio/kafka-go"

	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	aws "github.com/aws/aws-sdk-go-v2/config"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam_v2"
)

type config struct {
	Topic   string
	Brokers []string
	GroupID string
}

type Service interface {
	Shutdown(context.Context) error
}

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

	cfg, err := aws.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load aws config", "reason", err.Error())
	}

	sasl := &aws_msk_iam_v2.Mechanism{
		Signer:      signer.NewSigner(),
		Credentials: cfg.Credentials,
		Region:      c.Config().Topic,
		SignTime:    time.Now(),
		Expiry:      time.Minute * 15,
	}

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Topic:   c.Config().Topic,
		Brokers: c.Config().Brokers,
		GroupID: c.Config().GroupID,
		MaxWait: 50000 * time.Millisecond,
		Dialer: &kafka.Dialer{
			Timeout:       15 * time.Second,
			DualStack:     true,
			SASLMechanism: sasl,
			TLS: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	})

	go c.listen(ctx)

	return nil
}

func (c *consumer) Shutdown(ctx context.Context) error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("consumer service did not shutdown properly: %w", err)
	}
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
				logger.Error("failed to read message from kafka", "reason", err.Error())
				continue
			}
			logger.Info("received message", "offset", msg.Offset, "message", string(msg.Value))
		}
	}

	if err := c.reader.Close(); err != nil {
		logger.ErrorCtx(ctx, "failed to close kafka reader")
	}
}
