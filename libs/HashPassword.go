package lib

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashedPassword, password string) error {
	// Compare the hashed password with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// If passwords don't match, return an error
		return err
	}
	// Passwords match
	return nil
}
