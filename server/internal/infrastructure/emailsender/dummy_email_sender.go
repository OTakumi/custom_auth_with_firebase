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

// SendOTP simulates sending an OTP email.
// In development, check the Firestore Emulator UI to see the OTP.
func (s *DummyEmailSender) SendOTP(ctx context.Context, toEmail, otp string) error {
	log.Printf("Dummy Email Sent to: %s (check Firestore Emulator UI for OTP)", toEmail)

	return nil
}

// Ensure DummyEmailSender implements the EmailSender interface.
var _ emailsender.EmailSender = (*DummyEmailSender)(nil)
