package utils

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {
	// Check if the application is running in production
	environment := os.Getenv("ENVIRONMENT")
	// fmt.Println(os.Getenv("ENVIRONMENT"))
	if environment != "production" {
		// load .env file in non-production environments
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	return os.Getenv(key)
}

func CalculateAge(birthDate time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()
	// If this year's birthday has not occurr yet, subtract one from the age.
	if today.Month() < birthDate.Month() || (today.Month() == birthDate.Month() && today.Day() < birthDate.Day()) {
		age--
	}
	return age
}
