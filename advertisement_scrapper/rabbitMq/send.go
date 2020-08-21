package rabbitMq

import (
	"fmt"
	"os"

	"Krisha/advertisement_scrapper/ads"
	"Krisha/advertisement_scrapper/utils"

	"github.com/streadway/amqp"
)

func Receiver(){
	utils.Logging("Info","Starting the receiver..")

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
		"links", // name
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
			link := fmt.Sprintf("%s",d.Body)
			advertisement,err := ads.GetAdd(link)
			if err!= nil{
				d.Ack(false)
			}else {
				if  advertisement != ""{
					send(advertisement)
				}
			}
		}
	}()

	utils.Logging("Info"," [*] Waiting for messages.")
	<-forever
}

func send(add string){
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

	body := add
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	utils.Logging("Debug","Sent - "+body)
	if err!=nil{
		utils.Logging("Error","Failed to publish a message - "+err.Error())
		return
	}
}