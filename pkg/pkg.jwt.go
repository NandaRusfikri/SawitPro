package pkg

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Claims struct {
	jwt.StandardClaims
	Id          int    `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

func Sign(claims Claims) (string, error) {

	privateKeyPath := "cert/private_key.pem"

	// Baca private key dari file
	privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Println("Error reading private key:", err)
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return "", err
	}

	// Buat token menggunakan algoritma RSA-PSS (RS256)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign token dengan private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil

}

func StartsWithBearer(s string) bool {
	return len(s) > 7 && s[:7] == "Bearer "
}
func ExtractTokenFromHeader(authHeader string) string {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}
func GetPathCert(filename string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get working directory:", err)
		return "", err
	}
	filePath := filepath.Join(currentDir, "cert", filename)
	return filePath, nil
}

func Auth(ctx echo.Context) (*Claims, error) {

	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" || !StartsWithBearer(authHeader) {
		return nil, ctx.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	TokenString := ExtractTokenFromHeader(authHeader)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get working directory:", err)
		return nil, err
	}
	filePath := filepath.Join(currentDir, "cert", "public_key.pem")

	pubKey, err := loadPublicKey(filePath)
	if err != nil {
		log.Errorln("Failed to load public key:", err)
	}

	_, err = validateJWT(TokenString, pubKey)
	if err != nil {
		log.Errorln("Token validation failed:", err)
	}

	token, err := jwt.ParseWithClaims(TokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi token dengan menggunakan kunci publik
		return pubKey, nil
	})

	// Memeriksa apakah parsing berhasil atau terdapat error
	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	// Memeriksa apakah token valid dan dapat diassert ke dalam Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, nil

}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	//pubKeyBytes, err := ioutil.ReadFile(path)
	//if err != nil {
	//	return nil, err
	//}

	pubKeyBytes := []byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvj+gQB8yD/l/51cDO19L\nm4IJ20k3KXpYsRuZEfEMsDMYi4tMmQ/Iy2QxJUxq9+QjJGmabphMWY5EdtOsDtsK\nHFzyEjctvNErnoJclhaG8fEUo4EGLanft5gAIGcG08o3p29aJVhWajMo3drT7iUV\nBOlTMoArTuWyM4YmTGdRc528TCGxf1aQdvHTUZ1d13ungYUZPfw9KD82XhkUU44G\n+rbBSVkV889JmXQQo1xCdhifGQ4R/XHG6QQM8hBFl0ALajEKvCVPIOZwiBQ/nZhp\nfeD6GyHNtjA6xEH1ivoqnD+OJGQtwEa2MfIPJ72Xt9tOXVd7ImAYYX4NXQ/SLBZh\nrwIDAQAB\n-----END PUBLIC KEY-----\n")
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func validateJWT(tokenString string, pubKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Gunakan kunci publik untuk memverifikasi token
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
