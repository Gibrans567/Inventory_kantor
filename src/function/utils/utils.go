package utils

import "golang.org/x/crypto/bcrypt"

// CheckPasswordHash membandingkan password plaintext dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}