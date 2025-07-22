package services

import (
	"net/smtp"
	"os"
)

func SendEmail(to string, otp string) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"))
	message := []byte("To: " + to + "\r\n" +
		"Subject: OTP Verification\r\n" +
		"\r\n" +
		"Your OTP is: " + otp + "\r\n")
	return smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		[]string{to},
		message,
	)
}
