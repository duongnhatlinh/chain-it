package rabbitmq

import "github.com/streadway/amqp"

type Session struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func CreateSession(rabbitmqUrl string) Session {

	if rabbitmqUrl == "" {
		rabbitmqUrl = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitmqUrl)

	if err != nil {
		panic("Failed to connect to RabbitMQ" + err.Error())
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(err.Error())
	}

	return Session{
		Ch:   ch,
		Conn: conn,
	}
}

func (s Session) Close() {

	if s.Conn != nil {
		s.Conn.Close()
	}
	if s.Ch != nil {
		s.Ch.Close()
	}
}
