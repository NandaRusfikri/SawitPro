package util

import (
	"fmt"
	"regexp"
	"strings"
)

// Fungsi validasi untuk nomor telepon
func ValidatePhoneNumber(phoneNumber string) bool {
	// Cek apakah panjang nomor telepon sesuai dengan aturan
	if len(phoneNumber) < 11 || len(phoneNumber) > 14 {
		fmt.Println("len", len(phoneNumber))
		return false
	}

	// Cek apakah nomor telepon dimulai dengan "+62"
	if !strings.HasPrefix(phoneNumber, "+62") {
		return false
	}

	return true
}

// Fungsi validasi untuk nama lengkap
func ValidateFullName(fullName string) bool {
	// Cek apakah panjang nama lengkap sesuai dengan aturan
	if len(fullName) < 3 || len(fullName) > 60 {
		return false
	}
	return true
}

// Fungsi validasi untuk password
func ValidatePassword(password string) bool {
	// Cek apakah panjang password sesuai dengan aturan
	if len(password) < 6 || len(password) > 64 {
		return false
	}

	// Cek apakah password mengandung minimal 1 huruf besar, 1 angka, dan 1 karakter khusus
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		// Konversi char menjadi string karena regexp.MatchString mengharapkan string sebagai argumen
		charStr := string(char)

		if regexp.MustCompile("[A-Z]").MatchString(charStr) {
			hasUpper = true
		}
		if regexp.MustCompile("[0-9]").MatchString(charStr) {
			hasDigit = true
		}
		if regexp.MustCompile("[^A-Za-z0-9]").MatchString(charStr) {
			hasSpecial = true
		}
	}

	return hasUpper && hasDigit && hasSpecial
}
