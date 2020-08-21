package rabbitMq

import (
	"Krisha/links_scrapper/utils"
	"os"

	"github.com/streadway/amqp"
)

func Send(url string){
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

	body := url
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