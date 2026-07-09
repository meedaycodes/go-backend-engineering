// Package main is the entry point that demonstrates how interfaces decouple
// calling code from concrete implementations. sendNotification accepts any
// Sender — it doesn't know or care if it's email or SMS.
package main

import (
	"errors"
	"fmt"

	"github.com/meedaycodes/day02-types-interfaces/internal/notification"
)

// sendNotification accepts the Sender interface, not a concrete type.
// This means it works with EmailSender, SMSSender, or any future sender
// without modification. It uses errors.Is() to inspect wrapped errors
// and identify the root cause through the Unwrap chain.
func sendNotification(s notification.Sender, recipient string, message string) {

	err := s.Send(recipient, message)

	if err != nil {
		if errors.Is(err, notification.ErrEmptyMessage) {
			fmt.Printf("Error Message: %s\n", err)
		}

		if errors.Is(err, notification.ErrEmptyRecipient) {
			fmt.Printf("Error Message: %s\n", err)
		}

		return
	}

	fmt.Println("Message sent successfully")

}

// main creates concrete senders and passes them to sendNotification,
// which only sees the Sender interface. Tests both success and error paths.
func main() {
	email := notification.EmailSender{SMTPHost: "senderbender", Port: 6875}
	sms := notification.SMSSender{APIKey: "cata767w5bx", APISecret: "HYUUII7"}

	sendNotification(email, "habeebaramide@yahoo.com", "welcome cracked developer")
	sendNotification(sms, "fahma.a.g@gmail.com", "straight to 10k gbp per month")
	sendNotification(email, "", "whys is this empty")
	sendNotification(sms, "habeebaramide@yahoo.com", "")

}
