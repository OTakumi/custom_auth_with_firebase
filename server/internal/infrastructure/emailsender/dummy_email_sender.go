package emailsender

import (
	"context"
	"log"

	"custom_auth_api/internal/domain/emailsender"
)

// DummyEmailSender is a dummy implementation of the EmailSender interface that logs emails.
type DummyEmailSender struct{}

// NewDummyEmailSender creates a new DummyEmailSender.
func NewDummyEmailSender() *DummyEmailSender {
	return &DummyEmailSender{}
}

// SendOTP logs the OTP email content to the console.
func (s *DummyEmailSender) SendOTP(ctx context.Context, toEmail, otp string) error {
	log.Printf("Dummy Email Sent to: %s with OTP: %s", toEmail, otp)

	return nil
}

// Ensure DummyEmailSender implements the EmailSender interface.
var _ emailsender.EmailSender = (*DummyEmailSender)(nil)
