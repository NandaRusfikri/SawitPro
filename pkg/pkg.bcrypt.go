package pkg

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func generateSalt() string {
	// Ganti dengan algoritma random salt yang lebih aman jika digunakan secara production
	return "mysalt"
}

func HashPassword(password string) string {
	salt := generateSalt()
	pw := []byte(password + salt)
	result, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return string(result)
}

func ComparePassword(hashPassword string, password string) error {
	salt := generateSalt()
	pw := []byte(password + salt)
	hw := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(hw, pw)
	return err
}
