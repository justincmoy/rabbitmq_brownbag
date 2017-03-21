package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	sleep      = kingpin.Flag("sleep", "sleep between messages").Default("0").Int()
	port       = kingpin.Flag("port", "rabbitmq port").Default("5672").String()
	queueName  = kingpin.Flag("queue", "rabbitmq queue").Required().String()
	routingKey = kingpin.Flag("routing", "routing key").Default("").String()
)

func main() {
	kingpin.Parse()

	conn, _ := amqp.Dial("amqp://guest:guest@localhost:" + *port)
	defer conn.Close()

	channel, _ := conn.Channel()
	defer channel.Close()

	channel.Qos(1, 0, false)
	messages, _ := channel.Consume(
		*queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			fmt.Printf("%s\n", d.Body)
			time.Sleep(time.Duration(*sleep) * time.Second)
			d.Ack(false)
		}
	}()

	<-forever
}
