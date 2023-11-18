package profile

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type UserProfile struct {
	UserID        int
	URL           string
	Nickname      string
	ContactPublic bool
	Address       string
	Biography     string
	Organization  string
	Country       string
	SocialLinks   []string
}

const (
	username = "root"
	password = "12345"
	dbname   = "apigetway"
)

func getDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, dbname))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUserProfile(profile UserProfile) (int64, error) {
	db, err := getDBConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := "INSERT INTO user_profiles (UserID, URL, Nickname, ContactPublic, Address, Biography, Organization, Country, SocialLinks) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, profile.UserID, profile.URL, profile.Nickname, profile.ContactPublic, profile.Address, profile.Biography, profile.Organization, profile.Country, strings.Join(profile.SocialLinks, ","))
	if err != nil {
		return 0, err
	}

	userID, _ := result.LastInsertId()
	return userID, nil
}

func GetUserProfileByID(userID int) (UserProfile, error) {
	db, err := getDBConnection()
	if err != nil {
		return UserProfile{}, err
	}
	defer db.Close()

	var profile UserProfile
	query := "SELECT URL, Nickname, ContactPublic, Address, Biography, Organization, Country, SocialLinks FROM user_profiles WHERE UserID = ?"
	var socialLinks []byte // Cambio aquí para escanear SocialLinks como []byte

	err = db.QueryRow(query, userID).Scan(
		&profile.URL, &profile.Nickname, &profile.ContactPublic,
		&profile.Address, &profile.Biography, &profile.Organization,
		&profile.Country, &socialLinks,
	)
	if err != nil {
		return UserProfile{}, err
	}

	// Convertir el []byte de SocialLinks a string y luego a []string
	profile.SocialLinks = strings.Split(string(socialLinks), ",")

	return profile, nil
}

func UpdateUserProfile(userID int, profile UserProfile) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE user_profiles SET URL=?, Nickname=?, ContactPublic=?, Address=?, Biography=?, Organization=?, Country=?, SocialLinks=? WHERE UserID=?"

	_, err = db.Exec(query, profile.URL, profile.Nickname, profile.ContactPublic, profile.Address, profile.Biography, profile.Organization, profile.Country, strings.Join(profile.SocialLinks, ","), userID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserProfile(userID int) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM user_profiles WHERE UserID = ?"
	_, err = db.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

func createSampleUserProfiles() {
	profiles := []UserProfile{
		{
			URL:           "https://www.example.com/user1",
			Nickname:      "user1",
			ContactPublic: true,
			Address:       "123 Street, City",
			Biography:     "I'm user1",
			Organization:  "Org1",
			Country:       "Country1",
			SocialLinks:   []string{"https://twitter.com/user1", "https://github.com/user1"},
		},
		{
			URL:           "https://www.example.com/user2",
			Nickname:      "user2",
			ContactPublic: false,
			Address:       "456 Avenue, Town",
			Biography:     "I'm user2",
			Organization:  "Org2",
			Country:       "Country2",
			SocialLinks:   []string{"https://twitter.com/user2", "https://github.com/user2"},
		},
		{
			URL:           "https://www.example.com/user3",
			Nickname:      "user3",
			ContactPublic: true,
			Address:       "789 Road, Village",
			Biography:     "I'm user3",
			Organization:  "Org3",
			Country:       "Country3",
			SocialLinks:   []string{"https://twitter.com/user3", "https://github.com/user3"},
		},
	}

	for _, profile := range profiles {
		_, err := CreateUserProfile(profile)
		if err != nil {
			log.Printf("Error creando perfil: %s", err.Error())
		}
	}
}

func main1() {
	createSampleUserProfiles()
}
