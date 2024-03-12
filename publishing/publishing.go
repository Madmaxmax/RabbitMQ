package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Can't connect to RabbitMQ")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a channel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	rand.Seed(time.Now().UnixNano())

	addTask := struct {
		Number1 int `json:"number1"`
		Number2 int `json:"number2"`
	}{
		Number1: rand.Intn(999),
		Number2: rand.Intn(999),
	}

	body, err := json.Marshal(addTask)
	handleError(err, "Error encoding JSON")

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("AddTask: %d + %d", addTask.Number1, addTask.Number2)
}
