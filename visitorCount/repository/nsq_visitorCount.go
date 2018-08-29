package repository

import (
	nsq "github.com/nsqio/go-nsq"
	"github.com/virph/sc-project/visitorCount"
)

type nsqVisitorCountRepository struct {
	producer *nsq.Producer
}

func NewNsqVisitorCountRepository(producer *nsq.Producer) visitorCount.NsqVisitorCountRepository {
	return &nsqVisitorCountRepository{
		producer: producer,
	}
}

func (r *nsqVisitorCountRepository) Publish(topic string, payload []byte) error {
	return r.producer.Publish(topic, payload)
}

func (r *nsqVisitorCountRepository) AddVisitor() {
	r.Publish("visitor_count", []byte("wkwk"))
}
