package producer

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam_v2"
	"github.com/syllabix/kafkaless/reverser"

	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	aws "github.com/aws/aws-sdk-go-v2/config"
)

type config struct {
	Address []string
	Topic   string
}

type Service interface {
	EmitEvent(context.Context, string) error
}

type producer struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]
	reverser weaver.Ref[reverser.Service]
	kafka    *kafka.Writer
}

func (p *producer) Init(ctx context.Context) error {
	p.Logger().
		With("address", p.Config().Address).
		With("topic", p.Config().Topic).
		Info("initializing producer ...")

	cfg, err := aws.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load aws config", "reason", err.Error())
	}

	sasl := &aws_msk_iam_v2.Mechanism{
		Signer:      signer.NewSigner(),
		Credentials: cfg.Credentials,
		Region:      p.Config().Topic,
		SignTime:    time.Now(),
		Expiry:      time.Minute * 15,
	}

	p.kafka = &kafka.Writer{
		Addr:      kafka.TCP(p.Config().Address...),
		Topic:     p.Config().Topic,
		Balancer:  new(kafka.RoundRobin),
		BatchSize: 1,
		Transport: &kafka.Transport{
			SASL: sasl,
			TLS: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	return nil
}

func (p *producer) EmitEvent(ctx context.Context, str string) error {
	reversed, err := p.reverser.Get().Reverse(ctx, str)
	if err != nil {
		return fmt.Errorf("failed to emit event: %w", err)
	}

	err = p.kafka.WriteMessages(ctx, kafka.Message{
		Value: []byte(reversed),
	})
	if err != nil {
		p.Logger().Error("failed to emit event", "error", err)
		return fmt.Errorf("failed to emit event: %w", err)
	}

	return nil
}
