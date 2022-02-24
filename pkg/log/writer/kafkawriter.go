package writer

import (
	"github.com/Shopify/sarama"
)

type KafkaWriter struct {
	topic        string
	syncProducer sarama.SyncProducer
}

func NewKafkaWriter(topic string, syncProducer sarama.SyncProducer) *KafkaWriter {
	return &KafkaWriter{
		topic:        topic,
		syncProducer: syncProducer,
	}
}

func (self *KafkaWriter) Write(p []byte) (n int, err error) {
	topic := self.topic
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(p)
	//防止数据倾斜
	_, _, err = self.syncProducer.SendMessage(msg)
	return len(p), err
}
