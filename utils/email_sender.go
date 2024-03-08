package utils

import (
	"bytes"
	"fmt"
	"home_manager/config"
	"home_manager/entities"
	"net/smtp"
	"strings"
)

func SendVerificationEmail(email string, token entities.VerificationToken) error {
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
	message.WriteString("Please verify your email with this link\n")
	message.WriteString(baseUrl)
	message.WriteString("/verify")
	message.WriteString("?user_id=" + token.UserId)
	message.WriteString("&verify_token=" + token.Token)
	fmt.Println(message.String())

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf(message.String())))

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}
