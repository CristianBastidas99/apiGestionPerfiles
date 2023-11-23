package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CristianBastidas99/profile-service/profile"
	"github.com/streadway/amqp"
)

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro "userID" de la URL
	queryParams := r.URL.Query()
	userID := queryParams.Get("userID")

	// Convertir userID a entero
	profileID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedProfile profile.UserProfile
	err = json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Aquí llamamos a la función para actualizar el perfil con el userID
	err = profile.UpdateUserProfile(profileID, updatedProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Perfil actualizado exitosamente"))
}

func createProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newProfile profile.UserProfile
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Aquí llamamos a la función para crear un nuevo perfil
	// Supongamos que ya tienes acceso a la función de profile.go
	userID, err := profile.CreateUserProfile(newProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Perfil creado exitosamente con ID: %d", userID)))
}

func getProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro "userID" de la URL
	queryParams := r.URL.Query()
	userID := queryParams.Get("userID")

	// Convertir userID a entero
	profileID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Llamamos a la función para obtener el perfil del usuario por ID
	profile, err := profile.GetUserProfileByID(profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serializamos el perfil a JSON y respondemos con el perfil del usuario
	response, _ := json.Marshal(profile)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func deleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro "userID" de la URL
	queryParams := r.URL.Query()
	userID := queryParams.Get("userID")

	// Convertir userID a entero
	profileID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = profile.DeleteUserProfile(profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Perfil eliminado exitosamente"))
}

func createUserProfileFromRegistration(userID int) error {
	// Aquí puedes obtener información adicional del usuario si es necesario
	// por ejemplo, haciendo una llamada a la base de datos de usuarios registrados

	// Crea un perfil básico para el nuevo usuario
	newProfile := profile.UserProfile{
		UserID:        userID,
		URL:           "",
		Nickname:      "",
		ContactPublic: false,
		Address:       "",
		Biography:     "",
		Organization:  "",
		Country:       "",
		SocialLinks:   []string{},
	}

	// Crea el perfil para el usuario recién registrado
	_, err := profile.CreateUserProfile(newProfile)
	if err != nil {
		return err
	}

	return nil
}

func registrationWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registrationInfo struct {
		UserID int `json:"userID"`
		// Otros campos relevantes del registro, si los hay
	}

	err := json.NewDecoder(r.Body).Decode(&registrationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Crear un perfil para el nuevo usuario
	err = createUserProfileFromRegistration(registrationInfo.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Perfil creado para el nuevo usuario"))
}

func loggingMiddleware(next http.Handler, conn *amqp.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the request details for logging
		requestInfo := LogMessage{
			AppGener:    "profile_service",
			Tipo:        r.Method,
			ClaseModelo: r.URL.Path,
			FechaHora:   time.Now().Format("2006-01-02 15:04:05"),
			Resumen:     fmt.Sprintf("Received request: %s %s", r.Method, r.URL.Path),
		}

		// Proceed with the request
		next.ServeHTTP(w, r)

		// Create a response log message
		responseInfo := LogMessage{
			AppGener:    "profile_service",
			Tipo:        r.Method,
			ClaseModelo: r.URL.Path,
			FechaHora:   time.Now().Format("2006-01-02 15:04:05"),
			Resumen:     fmt.Sprintf("Processed request: %s %s", r.Method, r.URL.Path),
			Descripcion: fmt.Sprintf("Status code: %d" /*, w.StatusCode*/),
		}

		// Marshal the log messages to JSON format
		requestJSON, err := json.Marshal(requestInfo)
		if err != nil {
			log.Println("Error marshalling request log message:", err)
		}

		responseJSON, err := json.Marshal(responseInfo)
		if err != nil {
			log.Println("Error marshalling response log message:", err)
		}

		// Send the log messages to the RabbitMQ queue
		ch, err := conn.Channel()
		if err != nil {
			log.Println("Error opening channel:", err)
		}
		defer ch.Close()

		err = ch.Publish(
			"",
			"cola_1",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        requestJSON,
			},
		)
		if err != nil {
			log.Println("Error publishing request log message:", err)
		}

		err = ch.Publish(
			"",
			"cola_1",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        responseJSON,
			},
		)
		if err != nil {
			log.Println("Error publishing response log message:", err)
		}
	})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type LogMessage struct {
	AppGener    string
	Tipo        string
	ClaseModelo string
	FechaHora   string
	Resumen     string
	Descripcion string
}

func main() {

	router := http.NewServeMux()

	// Manejador para actualizar el perfil
	router.HandleFunc("/create-profile", createProfileHandler)
	router.HandleFunc("/update-profile", updateProfileHandler)
	router.HandleFunc("/get-profile", getProfileHandler)
	router.HandleFunc("/delete-profile", deleteProfileHandler)

	// Manejador para el webhook de registro de usuarios
	router.HandleFunc("/registration-webhook", registrationWebhookHandler)

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cola_1", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	failOnError(err, "Failed to declare a queue")

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
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	// Aplicar middleware para registrar invocaciones al servicio
	loggedRouter := loggingMiddleware(router, conn)

	log.Println("Servidor gestion de  iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
