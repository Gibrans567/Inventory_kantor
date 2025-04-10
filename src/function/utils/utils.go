package utils

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"time"
)

// Secret key untuk signing JWT
var jwtSecret = []byte("InventoryAppSecretKeyAwh1029")

// Fungsi untuk menghasilkan JWT token
func GenerateToken(userID uint, roleID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"roleID": roleID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Fungsi untuk memvalidasi token JWT
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}

// Fungsi untuk refresh token
func RefreshToken(tokenString string) (string, error) {
	// Memvalidasi token yang ada
	token, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Mendapatkan klaim dan menghasilkan token baru
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Menghasilkan token baru dengan klaim yang sama
	return GenerateToken(uint(claims["userID"].(float64)), uint(claims["roleID"].(float64)))
}
