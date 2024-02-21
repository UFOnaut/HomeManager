package utils

import (
	"fmt"
	"home_manager/config"
	"net/smtp"
	"strings"
)

func SendVerificationEmail(email string, verificationToken string) error {
	cfg := config.GetConfig()
	// Sender data.
	emailCredentials := cfg.EmailCredentials
	from := emailCredentials.Email
	password := emailCredentials.Password
	baseUrl := cfg.Endpoint.BaseUrl

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	var message strings.Builder
	message.WriteString("Please verify your email with this link ")
	message.WriteString(baseUrl)
	message.WriteString("?email=" + email)
	message.WriteString("&verify_token=" + verificationToken)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message.String()))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}
