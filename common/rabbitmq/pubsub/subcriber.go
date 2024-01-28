package pubsub

import (
	"encoding/json"

	"it-chain/common/rabbitmq"

	"github.com/DE-labtory/sdk/logger"
	"github.com/rs/xid"
)

type Message struct {
	MatchingValue string
	Data          []byte
}

type TopicSubscriber struct {
	rabbitmq.Session
	exchange string
}

func NewTopicSubscriber(rabbitmqUrl string, exchange string) *TopicSubscriber {

	session := rabbitmq.CreateSession(rabbitmqUrl)
	err := session.Ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		panic(err.Error())
	}

	return &TopicSubscriber{
		Session:  session,
		exchange: exchange,
	}
}

func (t *TopicSubscriber) SubscribeTopic(topic string, source interface{}) error {

	q, err := t.Session.Ch.QueueDeclare(
		xid.New().String(), // name
		false,              // durable
		true,               // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	err = t.Session.Ch.QueueBind(
		q.Name,     // queue name
		topic,      // routing key
		t.exchange, // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	msgs, err := t.Session.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	go func(queueName string) {

		p, _ := NewParamBasedRouter()
		p.SetHandler(q.Name, source)

		for delivery := range msgs {
			message := &Message{}
			data := delivery.Body

			if err := json.Unmarshal(data, message); err != nil {
				logger.Errorf(nil, "[Common] Fail to unmarshal rabbitmq message - Err: [%s]", err.Error())
			}

			if err := p.Route(queueName, message.Data, message.MatchingValue); err != nil {
				logger.Errorf(nil, "[Common] Fail to route rabbitmq message - Err: [%s]", err.Error())
			} //해당 event를 처리하기 위한 matching value 에는 structName이 들어간다.
		}
	}(q.Name)

	return nil
}

func (t *TopicSubscriber) Close() {
	t.Session.Close()
}
