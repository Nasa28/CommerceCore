package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(password []byte, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	fmt.Println(err)
	return err == nil
}
