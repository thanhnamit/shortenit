package internal

import (
	"context"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

type Consumer struct {
	config *config.Config
}

func InitKafkaConsumer(config *config.Config) {
	cgrpHandler := Consumer{config: config}
	handler := otelsarama.WrapConsumerGroupHandler(&cgrpHandler)

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	cgrp, err := sarama.NewConsumerGroup(config.BrokerList, "test_group", cfg)
	if err != nil {
		log.Fatalln("Failed to start consumer group: ", err)
	}

	err = cgrp.Consume(context.Background(), []string{config.GetUrlTopic}, handler)
	if err != nil {
		log.Fatalln("Failed to consume: ", err)
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		time.Sleep(4*time.Second)
		consumeMessage(message, c.config)
		session.MarkMessage(message, "")
	}
	return nil
}

func consumeMessage(message *sarama.ConsumerMessage, config *config.Config) {
	meter := otel.Meter("api-shortenit-v1")
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))

	tr := otel.Tracer(config.TracerName)
	_, span := tr.Start(ctx, "kafka.consumer.ConsumeMessage", trace.WithAttributes(semconv.MessagingOperationProcess))
	defer span.End()

	log.Println("Received message: ", string(message.Value))
	recordEventConsumed(meter, ctx, GetTag(config))
}

func recordEventConsumed(meter metric.Meter, ctx context.Context, l *label.KeyValue) {
	upDownCounter := metric.Must(meter).NewInt64UpDownCounter("api-shortenit-v1.kafka-queue-size")
	upDownCounter.Add(ctx, -1, *l)
}
