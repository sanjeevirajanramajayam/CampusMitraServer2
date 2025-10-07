package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}
func GenerateJWT(email, rollNo, role,user_name string) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"rollNo": rollNo,
		"role":   role,
		"user_name": user_name,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret()) // Read secret fresh
}
func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil // Read secret fresh
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return token.Claims.(jwt.MapClaims), nil
}
