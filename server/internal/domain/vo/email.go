package email

import (
	"errors"
	"regexp"
)

// Email represents an email address value object.
type Email struct {
	Value string
}

const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	emailPattern          = regexp.MustCompile(emailRegex)
)

// NewEmail creates a new Email value object.
func NewEmail(email string) (*Email, error) {
	if !emailPattern.MatchString(email) {
		return nil, ErrInvalidEmailFormat
	}

	return &Email{Value: email}, nil
}

// Verify if it has become the general format of an email.
