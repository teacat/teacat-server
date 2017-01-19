package auth

import "golang.org/x/crypto/bcrypt"

func Encrypt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
