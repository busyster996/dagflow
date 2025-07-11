package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/busyster996/dagflow/pkg/rabbitmq"
)

func main() {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsConsumerName("consumer"),
	)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("awaiting signal")
		sig := <-sigs
		fmt.Println(sig)
		fmt.Println("stopping consumer")

		consumer.Close()
	}()

	// block main thread - wait for shutdown signal
	err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("consumed: %v", string(d.Body))
		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	})
	if err != nil {
		log.Fatal(err)
	}
}
