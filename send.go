package main

import (
	"strconv"

	"github.com/streadway/amqp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	body         = kingpin.Flag("body", "body of message").String()
	exchangeName = kingpin.Flag("exchange", "rabbitmq exchange").Default("").String()
	number       = kingpin.Flag("number", "number of messages to send").Default("1").Int()
	port         = kingpin.Flag("port", "rabbitmq port").Default("5672").String()
	queueName    = kingpin.Flag("queue", "rabbitmq queue").Default("").String()
	routingKey   = kingpin.Flag("routing", "routing key").Default("").String()
)

func main() {
	kingpin.Parse()

	if *queueName == "" && *exchangeName == "" {
		kingpin.FatalUsage("Must define --queue or --exchange")
	}
	if *queueName != "" && *exchangeName != "" {
		kingpin.FatalUsage("Must define --queue or --exchange")
	}
	if *queueName != "" {
		*routingKey = *queueName
		*exchangeName = ""
	}

	conn, _ := amqp.Dial("amqp://guest:guest@localhost:" + *port)
	defer conn.Close()

	channel, _ := conn.Channel()
	defer channel.Close()

	for i := 0; i < *number; i++ {
		var bodyString = *body
		if bodyString == "" {
			bodyString = strconv.Itoa(i)
		}

		channel.Publish(
			*exchangeName,
			*routingKey,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				Body:        []byte(bodyString),
				ContentType: "text/plain",
			})
	}
}
