package internal

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
	cfg      *config.Config
}

func (ep KafkaProducer) Publish(ctx context.Context, event core.GetUrlEvent, topic string) {
	tr := otel.Tracer(ep.cfg.TracerName)
	_, span := tr.Start(ctx, "kafka.producer.Publish")
	defer span.End()

	json, _ := json.Marshal(event)

	log.Println("Publishing event to topic: ", ep.cfg.GetUrlTopic)

	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(json),
	}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(&msg))

	ep.producer.Input() <- &msg
	successMsg := <-ep.producer.Successes()
	log.Println("Published with offset: ", successMsg.Offset)

	err := ep.producer.Close()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Fatalln("Error closing producer: ", err)
	}
}

func NewKafkaProducer(config *config.Config) KafkaProducer {
	brokerList := config.BrokerList
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokerList, cfg)
	if err != nil {
		log.Fatalln("Failed to start Sarama client: ", err)
	}

	producer = otelsarama.WrapAsyncProducer(cfg, producer)

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write message:", err)
		}
	}()

	return KafkaProducer{
		producer: producer,
		cfg:      config,
	}
}
