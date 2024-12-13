package helpers

import (
	"errors"
	"example/project-management-system/internal/models"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Helper function to validate user input
func ValidateUser(user *models.User) error {
    if user.Username == "" {
        return errors.New("username cannot be empty")
    }
    if user.Email == "" || !IsValidEmail(user.Email) {
        return errors.New("invalid email address")
    }
    if len(user.Password) < 8 {
        return errors.New("password must be at least 8 characters long")
    }
    return nil
}

// Helper function to check for duplicate key errors
func IsDuplicateKeyError(err error) bool {
    // This is a MySQL/PostgreSQL specific check
    // You might need to adjust for your specific database
    return strings.Contains(err.Error(), "duplicate key") || 
        strings.Contains(err.Error(), "unique constraint")
}

// Custom error for existing user
var ErrUserAlreadyExists = errors.New("user already exists")

// Placeholder for email validation 
func IsValidEmail(email string) bool {
    // Implement email validation logic
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}

// Placeholder for password hashing
func HashPassword(password string) (string, error) {
    // Use a secure hashing method like bcrypt
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}