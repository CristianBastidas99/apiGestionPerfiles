package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CristianBastidas99/profile-service/profile"
)

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
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

func CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
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

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
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

func CreateUserProfileFromRegistration(userID int) error {
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

func RegistrationWebhookHandler(w http.ResponseWriter, r *http.Request) {
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
	err = CreateUserProfileFromRegistration(registrationInfo.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Perfil creado para el nuevo usuario"))
}
