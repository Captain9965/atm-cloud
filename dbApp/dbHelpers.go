package dbApp

import(
	"golang.org/x/crypto/bcrypt" // Import for password hashing
)

// HashPassword function to hash passwords using bcrypt
func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}