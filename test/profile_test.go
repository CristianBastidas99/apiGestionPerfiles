package test

import (
	"testing"

	"github.com/CristianBastidas99/profile-service/profile"
)

func TestCreateUserProfile(t *testing.T) {
	newProfile := profile.UserProfile{
		URL:           "https://example.com",
		Nickname:      "testuser",
		ContactPublic: true,
		Address:       "123 Test St",
		Biography:     "Lorem ipsum",
		Organization:  "TestOrg",
		Country:       "TestCountry",
		SocialLinks:   []string{"link1", "link2"},
	}

	userID, err := profile.CreateUserProfile(newProfile)
	if err != nil {
		t.Errorf("Error creando perfil: %v", err)
	}

	// Comprobación de que se haya creado un ID válido
	if userID == 0 {
		t.Errorf("Se esperaba un ID de usuario válido, se obtuvo 0")
	}

	// También puedes realizar otras aserciones aquí, como verificar los valores insertados en la base de datos
}

func TestGetUserProfileByID(t *testing.T) {
	// Supongamos que tenemos un ID de usuario existente para probar
	userID := 4

	profile, err := profile.GetUserProfileByID(userID)
	if err != nil {
		t.Errorf("Error obteniendo perfil: %v", err)
	}

	// Realizar aserciones sobre los valores del perfil recuperado
	if profile.URL != "https://example.com" {
		t.Errorf("URL incorrecta en el perfil obtenido")
	}
	// Verificar otros campos del perfil aquí...
}

func TestUpdateUserProfile(t *testing.T) {
	updatedProfile := profile.UserProfile{
		URL:           "https://updated.com",
		Nickname:      "updateduser",
		ContactPublic: false,
		Address:       "789 Test St",
		Biography:     "Updated bio",
		Organization:  "UpdatedOrg",
		Country:       "UpdatedCountry",
		SocialLinks:   []string{"link5", "link6"},
	}

	userID := 3 // Suponiendo que el ID 1 existe para actualizar
	err := profile.UpdateUserProfile(userID, updatedProfile)
	if err != nil {
		t.Errorf("Error actualizando perfil: %v", err)
	}

	// Verificar la actualización obteniendo el perfil y comparando los valores actualizados
}

func TestDeleteUserProfile(t *testing.T) {
	userID := 3 // Suponiendo que el ID 1 existe para eliminar

	err := profile.DeleteUserProfile(userID)
	if err != nil {
		t.Errorf("Error eliminando perfil: %v", err)
	}

	// Intentar obtener el perfil eliminado para verificar si existe
	profile, err := profile.GetUserProfileByID(userID)
	if err == nil {
		t.Errorf("Se esperaba un error al obtener el perfil eliminado, se obtuvo un perfil: %v", profile)
	}
}
