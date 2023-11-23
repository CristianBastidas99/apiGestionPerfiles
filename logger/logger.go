package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

type LogMessage struct {
	AppGener    string
	Tipo        string
	ClaseModelo string
	FechaHora   string
	Resumen     string
	Descripcion string
}

var (
	AmqpURI = "amqp://guest:guest@rabbitmq:5672/" // Cambia por la URL de conexión de tu servidor RabbitMQ
	Queue   = "cola_1"                            // Nombre de la cola
)

func LoggingMiddleware(next http.Handler, conn *amqp.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the request details for logging
		requestInfo := LogMessage{
			AppGener:    "profile_service",
			Tipo:        r.Method,
			ClaseModelo: r.URL.Path,
			FechaHora:   time.Now().Format("2006-01-02 15:04:05"),
			Resumen:     fmt.Sprintf("Received request: %s %s", r.Method, r.URL.Path),
		}

		rec := NewCapturingResponseWriter(w)
		// Proceed with the request
		next.ServeHTTP(rec, r)

		// Create a response log message
		responseInfo := LogMessage{
			AppGener:    "profile_service",
			Tipo:        r.Method,
			ClaseModelo: r.URL.Path,
			FechaHora:   time.Now().Format("2006-01-02 15:04:05"),
			Resumen:     fmt.Sprintf("Processed request: %s %s", r.Method, r.URL.Path),
			Descripcion: fmt.Sprintf("Status code: %s", http.StatusText(rec.Status())),
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
			Queue,
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
			Queue,
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

// NewCapturingResponseWriter crea un ResponseWriter capturador para rastrear el código de estado de la respuesta.
func NewCapturingResponseWriter(w http.ResponseWriter) *CapturingResponseWriter {
	return &CapturingResponseWriter{w, http.StatusOK}
}

// CapturingResponseWriter implementa http.ResponseWriter para capturar el código de estado de la respuesta.
type CapturingResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader registra el código de estado de la respuesta.
func (w *CapturingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Status devuelve el código de estado capturado.
func (w *CapturingResponseWriter) Status() int {
	return w.status
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
