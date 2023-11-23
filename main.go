package main

import (
	"log"
	"net/http"

	"github.com/CristianBastidas99/profile-service/handlers"
	"github.com/CristianBastidas99/profile-service/logger"
	"github.com/streadway/amqp"
)

func main() {

	router := http.NewServeMux()

	// Manejador para actualizar el perfil
	router.HandleFunc("/create-profile", handlers.CreateProfileHandler)
	router.HandleFunc("/update-profile", handlers.UpdateProfileHandler)
	router.HandleFunc("/get-profile", handlers.GetProfileHandler)
	router.HandleFunc("/delete-profile", handlers.DeleteProfileHandler)

	// Manejador para el webhook de registro de usuarios
	router.HandleFunc("/registration-webhook", handlers.RegistrationWebhookHandler)

	conn, err := amqp.Dial(logger.AmqpURI)
	logger.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logger.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		logger.Queue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	logger.FailOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	logger.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	// Aplicar middleware para registrar invocaciones al servicio
	loggedRouter := logger.LoggingMiddleware(router, conn)

	log.Println("Servidor gestion de  iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
