package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func errHelper(err error){
	if err != nil {
		panic(err)
	}
}

func initRabbitMQ()*amqp.Connection{
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	errHelper(err)
	return conn
}

func main(){
	conn := initRabbitMQ()
	defer conn.Close()
	ch,err := conn.Channel()
	errHelper(err)
	defer ch.Close()
     
    producer(ch)  // To send messages
    // consumer(ch)     // To receive messages


}