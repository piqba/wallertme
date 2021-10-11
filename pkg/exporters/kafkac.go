package exporters

import "github.com/confluentinc/confluent-kafka-go/kafka"

const (
	// TXS_TOPIC_KEY ...
	TXS_TOPIC_KEY = "txs"
)

var (
	ProducerKafka = GetProducerClientKafka()
)

func GetProducerClientKafka() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	// defer p.Close()
	return p
}
