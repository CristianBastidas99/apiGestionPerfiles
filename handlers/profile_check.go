package handlers

import "net/http"

// HealthCheckHandler verifica el estado general del servicio
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Estado general: UP"))
}

// ReadyCheckHandler verifica si el servicio está listo
func ReadyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Estado del servicio: READY"))
}

// LiveCheckHandler verifica si el servicio está en vivo
func LiveCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Estado del servicio: ALIVE"))
}
