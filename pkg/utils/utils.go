package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)


func GetEnvVariable(key string) string{
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Unable to load environmental variable ")
	}
	return os.Getenv(key)
}


func HashPassword(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPwd), err

}