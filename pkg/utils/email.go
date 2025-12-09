package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendOTP(toEmail string, otp string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "MMGRAPP <"+os.Getenv("SENDER_EMAIL")+">")
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Your Money Manager Apps OTP Code")

	body := fmt.Sprintf(`
			Your OTP code is: %s

			This code is valid for 5 minutes.
			If you did not request this, please ignore this email.
			`, otp)

	mailer.SetBody("text/plain", body)

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return dialer.DialAndSend(mailer)
}
