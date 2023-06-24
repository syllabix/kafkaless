package producer

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
	"github.com/segmentio/kafka-go"
	"github.com/syllabix/kafkaless/reverser"
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

	p.kafka = &kafka.Writer{
		Addr:      kafka.TCP(p.Config().Address...),
		Topic:     p.Config().Topic,
		Balancer:  new(kafka.RoundRobin),
		BatchSize: 1,
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
