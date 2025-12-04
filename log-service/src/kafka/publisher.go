package kafka

import (
	"log"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/IBM/sarama"
)

type LogProducer interface {
	ConnectLogProducer([]string) (sarama.SyncProducer, error)
	PushLogsToQueue([]byte) error
}

type logs struct {
	conf *config.Config
}

func NewLogProducer(cfg *config.Config) LogProducer {
	return &logs{conf: cfg}
}

func (l *logs) ConnectLogProducer(brokers []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.RequiredAcks(l.conf.Kafka.RequiredAcks)
	config.Producer.Retry.Max = l.conf.Kafka.MaxRetry

	return sarama.NewSyncProducer(brokers, config)
}

func (l *logs) PushLogsToQueue(message []byte) error {

	brokers := l.conf.Kafka.Brokers

	prod, err := l.ConnectLogProducer(brokers)

	if err != nil {
		return err
	}

	defer prod.Close()

	msg := &sarama.ProducerMessage{
		Topic: l.conf.Kafka.Topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := prod.SendMessage(msg)

	if err != nil {
		return err
	}
	log.Printf("Log in published in %v %v", partition, offset)
	return nil
}
