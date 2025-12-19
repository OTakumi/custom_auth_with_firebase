package emailsender

import "context"

// EmailSender defines the interface for sending emails.
type EmailSender interface {
	SendOTP(ctx context.Context, toEmail, otp string) error
}
