// Package notification provides types and interfaces for sending notifications
// across different channels (email, SMS, etc.). It demonstrates Go's interface
// system — concrete senders implement the Sender interface independently,
// allowing calling code to work with any sender without knowing its type.
package notification

import (
	"errors"
	"fmt"
)

// Sentinel errors for validation failures. These are package-level variables
// (not constants) because errors.New returns a pointer, which is a runtime
// operation. Callers use errors.Is() to match against these after unwrapping.
var (
	ErrEmptyRecipient = errors.New("recipient cannot be empty")
	ErrEmptyMessage   = errors.New("message cannot be empty")
)

// Sender defines the contract for any notification channel.
// Any type that implements Send with this exact signature satisfies
// the interface implicitly — no "implements" keyword needed.
// This allows sendNotification in main.go to accept any sender
// without knowing whether it's email, SMS, or something else.
type Sender interface {
	Send(recipient string, message string) error
}

// EmailSender is a concrete implementation of Sender that sends via email.
// Fields are exported (capitalized) so they can be set from outside this package.
// SMTPHost is string because hostnames are text. Port is int because port
// numbers are numeric values with a valid range (1-65535).
type EmailSender struct {
	SMTPHost string
	Port     int
}

// SMSSender is a concrete implementation of Sender that sends via SMS.
// APIKey and APISecret are strings because API credentials are opaque
// text tokens — stored and passed, never computed on.
type SMSSender struct {
	APIKey    string
	APISecret string
}

// NotificationError is a custom error type that wraps a sentinel error
// with additional context — which sender type failed and for which recipient.
// It satisfies the error interface via Error() and supports unwrapping
// via Unwrap() so callers can use errors.Is() to inspect the root cause.
type NotificationError struct {
	SenderType string
	Recipient  string
	Err        error
}

// Error returns a formatted string combining the underlying error, sender type,
// and recipient. This method makes *NotificationError satisfy the error interface.
// Pointer receiver (*NotificationError) is required because Send methods return
// &NotificationError (a pointer), and the interface must be satisfied on the pointer type.
func (n *NotificationError) Error() string {
	return fmt.Sprintf("%s error from sendertype %s to recipient %s", n.Err, n.SenderType, n.Recipient)
}

// Unwrap returns the underlying sentinel error (e.g. ErrEmptyRecipient).
// This enables errors.Is() to see through the NotificationError wrapper
// and match against the original sentinel error.
func (n *NotificationError) Unwrap() error {
	return n.Err
}

// Send validates inputs and simulates sending an email. Returns a
// *NotificationError wrapping the sentinel error on failure, nil on success.
// Value receiver (EmailSender, not *EmailSender) because Send only reads
// fields — it doesn't modify the struct.
func (e EmailSender) Send(recipient string, message string) error {
	if recipient == "" {
		return &NotificationError{SenderType: "email", Recipient: recipient, Err: ErrEmptyRecipient}
	}
	if message == "" {
		return &NotificationError{SenderType: "email", Recipient: recipient, Err: ErrEmptyMessage}
	}

	fmt.Printf("[EMAIL] To: %s via %s:%d\n", recipient, e.SMTPHost, e.Port)
	fmt.Printf("[EMAIL] Body: %s\n", message)

	return nil
}

// Send validates inputs and simulates sending an SMS. Same pattern as
// EmailSender.Send — both independently satisfy the Sender interface
// without knowing about each other.
func (s SMSSender) Send(recipient string, message string) error {
	if recipient == "" {
		return &NotificationError{SenderType: "sms", Recipient: recipient, Err: ErrEmptyRecipient}
	}

	if message == "" {
		return &NotificationError{SenderType: "sms", Recipient: recipient, Err: ErrEmptyMessage}
	}

	fmt.Printf("[SMS] To: %s\n", recipient)
	fmt.Printf("[SMS] Body: %s\n", message)

	return nil
}
