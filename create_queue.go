package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	defaultUuid = uuid.NewV4().String()

	exchangeName = kingpin.Flag("exchange", "exchange name").Default("").String()
	exchangeType = kingpin.Flag("exchangetype", "exchange type").Default("fanout").String()
	port         = kingpin.Flag("port", "rabbitmq port").Default("5672").String()
	queueName    = kingpin.Flag("queue", "rabbitmq queue").Default(defaultUuid).String()
	routingKey   = kingpin.Flag("routing", "routing key").Default("").String()
)

func main() {
	kingpin.Parse()

	conn, _ := amqp.Dial("amqp://guest:guest@localhost:" + *port)
	defer conn.Close()

	channel, _ := conn.Channel()
	defer channel.Close()

	queue, _ := channel.QueueDeclare(
		*queueName,
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	fmt.Println("Created queue with name " + queue.Name)

	if *exchangeName != "" {
		channel.ExchangeDeclare(
			*exchangeName,
			*exchangeType,
			true,  // durable
			false, // auto-delete
			false, // internal
			false, // no-wait
			nil,   // arguments
		)

		channel.QueueBind(
			queue.Name,
			*routingKey,
			*exchangeName,
			false, // no-wait
			nil,   // arguments
		)

		fmt.Println("Binding " + *exchangeType + " exchange " + *exchangeName + " to queue " + queue.Name)
	}
}
