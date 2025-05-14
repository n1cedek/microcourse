package utils

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashedPass, reqPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(reqPass))
	return err == nil
}
