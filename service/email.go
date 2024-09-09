package service

import (
	"auth_service/storage/cache"
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	SMTPServer  string
	SMTPPort    string
	Username    string
	Password    string
	SenderEmail string
	cache       *cache.EmailCache
}

func NewEmailSender(smtpServer, smtpPort, username, password, senderEmail string) *EmailSender {
	return &EmailSender{
		SMTPServer:  smtpServer, //"smtp.gmail.com"
		SMTPPort:    smtpPort,   //"587"
		Username:    username,
		Password:    password,    //"xsay zgvy uuvd xven"
		SenderEmail: senderEmail, //"abdusamatovjavohir@gmail.com"
	}
}

func (e *EmailSender) SendVerificationEmail(toEmail, verificationLink string) error {
	e.cache.SaveLink(toEmail, verificationLink)

	subject := "Subject: Please Verify Your Email Address\n"
	body := fmt.Sprintf("Hello,\n\nPlease verify your email address by clicking the following link:\n%s\n\nThank you!", verificationLink)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPServer)

	err := smtp.SendMail(
		e.SMTPServer+":"+e.SMTPPort,
		auth,
		e.SenderEmail,
		[]string{toEmail},
		message,
	)

	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}
