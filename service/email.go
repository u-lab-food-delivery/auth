package service

import (
	"auth_service/config"
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

func NewEmailSender(cnf config.EmailSenderConfig, emailCacher *cache.EmailCache) *EmailSender {
	return &EmailSender{
		SMTPServer:  cnf.SMTPServer, //"smtp.gmail.com"
		SMTPPort:    cnf.SMTPPort,   //"587"
		Username:    "email",
		Password:    cnf.Password,
		SenderEmail: cnf.SenderEmail, //"abdusamatovjavohir@gmail.com"
		cache:       emailCacher,
	}
}

func (e *EmailSender) SendVerificationEmail(toEmail, verificationLink string) error {
	err := e.cache.SaveLink(toEmail, verificationLink)
	if err != nil {
		return err
	}

	subject := "Subject: Please Verify Your Email Address\n"
	body := fmt.Sprintf("Hello,\n\nPlease verify your email address by clicking the following link:\n%s\n\nThank you!", verificationLink)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", e.SenderEmail, e.Password, e.SMTPServer)

	err = smtp.SendMail(
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
