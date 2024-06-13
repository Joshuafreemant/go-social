package helpers

import (
	"github.com/Joshuafreemant/go-social/config"
	"gopkg.in/gomail.v2"
)

// SendEmail sends an email with the OTP
func SendEmail(email, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Config("EMAIL"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/plain", "Your OTP code is: "+otp)

	d := gomail.NewDialer("smtp.gmail.com", 465, config.Config("EMAIL"), config.Config("EMAIL_PASSWORD"))

	return d.DialAndSend(m)
}
