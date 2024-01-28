package pubsub

import (
	"encoding/json"
	"errors"
	"reflect"

	"it-chain/common/rabbitmq"

	"github.com/streadway/amqp"
)

type TopicPublisher struct {
	rabbitmq.Session
	exchange string
}

func NewTopicPublisher(rabbitmqUrl string, exchange string) *TopicPublisher {

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

	return &TopicPublisher{
		Session:  session,
		exchange: exchange,
	}
}

func (t *TopicPublisher) Publish(topic string, data interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	err = t.Ch.ExchangeDeclare(
		t.exchange, // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)

	if err != nil {
		return err
	}

	var matchingValue string

	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		matchingValue = reflect.TypeOf(data).Elem().Name()
	} else {
		matchingValue = reflect.TypeOf(data).Name()
	}

	b, err := json.Marshal(data)

	message := Message{
		MatchingValue: matchingValue,
		Data:          b,
	}

	byte, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err != nil {
		panic("Failed to open exchange" + err.Error())
	}

	err = t.Ch.Publish(
		t.exchange, // exchange
		topic,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byte,
		})

	if err != nil {
		return err
	}

	return nil
}

func (t *TopicPublisher) Close() {
	t.Session.Close()
}
