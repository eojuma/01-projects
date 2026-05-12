package utils

import (
	"errors"
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	TempPassword = "Password123"
	TempUsername = "forumuser"
	TempEmail    = "forumuser@gmail.com"
)

func ValidatePassword(password string) error {

	// Check if the password is empty as well as its length
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	} else if len(password) > 64 {
		return errors.New("password must be no more than 64 characters long")
	}

	// requires numbers
	var hasNumber bool
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	// requires uppercase letters
	var hasUpper bool
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	// requires lowercase letters
	var hasLower bool
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	// requires reject non printable characters
	for _, char := range password {
		if !unicode.IsPrint(char) {
			return errors.New("password must not contain non-printable characters")
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	// must not be empty
	if len(email) == 0 {
		return errors.New("email must not be empty")
	}

	// contains '@' symbol
	if !strings.Contains(email, "@") {
		return errors.New("email must contain '@' symbol")
	}

	// must not contain spaces
	if strings.Contains(email, " ") {
		return errors.New("email must not contain spaces")
	}

	// check if email domain is valid (eg. gmail.com or yahoo.com etc.)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("email must contain exactly one '@' symbol")
	}
	domain := parts[1]
	if mxContains, err := net.LookupMX(domain); err != nil {
		return errors.New("Error validating email domain: " + err.Error())
	} else if len(mxContains) == 0 {
		return errors.New("Error validating email domain: no MX records found for domain")
	}

	return nil
}

func ValidateUsername(username string) error {
	// must not be empty and must be between 3 and 32 characters long
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 32 {
		return errors.New("username must be no more than 32 characters long")
	}

	// must not contain spaces
	if strings.Contains(username, " ") {
		return errors.New("username must not contain spaces")
	}

	return nil
}

func VerifyEmail(email string) error {
	if email != TempEmail {
		return errors.New("you entered an invalid email")
	}
	return nil
}

func VerifyPassword(password string) error {
	if password != TempPassword {
		return errors.New("your password is incorrect")
	}
	return nil
}

func VerifyUsername(username string) error {
	if username != TempUsername {
		return errors.New("you ented an invalid username, please check again")
	}
	return nil
}

func GenerateUUID() string {
	newUUID, _ := uuid.NewV4()

	return newUUID.String()
}

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


func HashPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "Error", err
	}
	return string(hash), nil
}

func ComparePass(password, hash string) bool {
	hashed := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return hashed == nil
}
