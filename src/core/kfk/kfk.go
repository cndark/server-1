package kfk

import (
	"encoding/json"
	"fw/src/core"
	"fw/src/core/log"

	"github.com/Shopify/sarama"
)

// ============================================================================

type Kafka struct {
	c sarama.Client
	p sarama.SyncProducer
}

// ============================================================================

func NewKafka() *Kafka {
	return &Kafka{}
}

// ============================================================================

func (self *Kafka) Open(addrs []string) {
	conf := sarama.NewConfig()
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true // sync-producer needs true
	conf.Producer.Return.Errors = true    // sync-producer needs true

	c, err := sarama.NewClient(addrs, conf)
	if err != nil {
		core.Panic("open kafka client failed:", err)
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		core.Panic("open kafka producer failed:", err)
	}

	self.c = c
	self.p = p
}

func (self *Kafka) Close() {
	if self.p != nil {
		self.p.Close()
	}

	if self.c != nil {
		self.c.Close()
	}
}

func (self *Kafka) Producer(topic string, obj interface{}) {
	msg, err := json.Marshal(obj)
	if err != nil {
		log.Error("kafka producer json marshal err:", err)
		log.Error(core.Callstack())
		return
	}

	pid, offset, err := self.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	})

	if err != nil {
		log.Error("kafka producer send error:", err, pid, offset)
		log.Error(core.Callstack())
		return
	}
}
