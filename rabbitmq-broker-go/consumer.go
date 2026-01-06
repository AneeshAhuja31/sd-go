package main

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)


func consume(ch *amqp.Channel,q*amqp.Queue,ctx context.Context){
	msgs,err := ch.ConsumeWithContext(
		ctx,
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	errHelper(err)
	for msg := range msgs {
		var order RestaurantOrder
		err := json.Unmarshal(msg.Body,&order)
		errHelper(err)
        fmt.Printf("Received Order: ID=%d, Items=%v\n", order.Id, order.Items)
	}
}

func consumer(ch *amqp.Channel){
	q := initQueue(ch)
	ctx := context.Background()
	fmt.Println("Consumer started. Waiting for messages...")
    consume(ch, &q, ctx)
}
