package auth

import "golang.org/x/crypto/bcrypt"

func Encrypt(source *string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*source), bcrypt.DefaultCost)
	hashedString := string(hashedBytes)
	source = &hashedString
	return err
}

func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
