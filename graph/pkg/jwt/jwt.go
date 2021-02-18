package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SecretKey to sign tokens
var (
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt token and assign an username
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatalf("Error to generate key: %v", err)
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}
	return "", err
}
