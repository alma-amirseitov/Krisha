package rabbitMq

import (
	"fmt"
	"os"

	"Krisha/saver/elastic"
	"Krisha/saver/utils"


	"github.com/streadway/amqp"
)
func Receiver(){
	utils.Logging("Info","Starting the saver receiver..")

	host := os.Getenv("KRISHA_RABBIT")
	conn, err := amqp.Dial(host)

	if err!=nil{
		utils.Logging("Error","Failed to connect to RabbitMQ - "+err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err!=nil{
		utils.Logging("Error","Failed to open a channel - "+err.Error())
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"advertisement", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err!=nil{
		utils.Logging("Error","Failed to declare a queue - "+err.Error())
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue,
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err!=nil{
		utils.Logging("Error","Failed to register a consumer - "+err.Error())
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			advertisement := fmt.Sprintf("%s",d.Body)
			elastic.ToElastic(advertisement)
		}
	}()

	fmt.Println("Info"," [*] Waiting for messages.")
	<-forever
}