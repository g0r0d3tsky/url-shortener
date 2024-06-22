package kafka

import (
	"encoding/json"
	"log/slog"
	"time"
	"url-service/config"
	"url-service/internal/domain"

	"github.com/IBM/sarama"
)

type Producer struct {
	sarama.AsyncProducer
}

func New(config *config.Config) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.DefaultVersion
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(config.Kafka.BrokerList, cfg)
	if err != nil {
		slog.Error("create producer", err)
	}
	go func() {
		for err = range producer.Errors() {
			slog.Error("write access log entry", err)
		}
	}()

	return &Producer{producer}, nil
}

func (p Producer) Push(topic string, url *domain.Url) error {
	jsonMessage, err := json.Marshal(url)
	if err != nil {
		slog.Error("marshal message:", err)
		return err
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}
	p.Input() <- kafkaMessage

	return nil
}
