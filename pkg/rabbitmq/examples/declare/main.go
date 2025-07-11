package main

import (
	"log"

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

	declare, err := rabbitmq.NewDeclarator(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer declare.Close()

	// declare dead-letter queue
	err = declare.Queue(rabbitmq.QueueOptions{
		Name:    "re_my_queue",
		Declare: true,
		Args: map[string]interface{}{
			"x-dead-letter-exchange":    "events",
			"x-dead-letter-routing-key": "my_routing_key",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	err = declare.BindQueues([]rabbitmq.Binding{
		{
			Source:      "events",
			RoutingKey:  "my_routing_key",
			Destination: "re_my_queue",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
