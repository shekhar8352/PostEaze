package email_service

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/shekhar8352/PostEaze/constants"
)

type EmailService interface {
	SendEmail(to []string, subject string, body string) error
	SendNotificationEmail(to string, notificationContent string) error
	SendTeamInviteEmail(to string, inviteLink string) error
}

type GmailEmailService struct {
	smtpHost string
	smtpPort string
	email    string
	password string
}

func NewGmailEmailService() *GmailEmailService {
	return &GmailEmailService{
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
		email:    os.Getenv("SMTP_EMAIL"),
		password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *GmailEmailService) SendEmail(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.email, s.password, s.smtpHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte("Subject: " + subject + "\n" + mime + "\n" + body)

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	if err := smtp.SendMail(addr, auth, s.email, to, msg); err != nil {
		return err
	}
	return nil
}

func (s *GmailEmailService) SendNotificationEmail(to string, notificationContent string) error {
	body := fmt.Sprintf(constants.NotificationEmailBody, notificationContent)
	return s.SendEmail([]string{to}, constants.NotificationEmailSubject, body)
}

func (s *GmailEmailService) SendTeamInviteEmail(to string, inviteLink string) error {
	body := fmt.Sprintf(constants.TeamInviteEmailBody, inviteLink)
	return s.SendEmail([]string{to}, constants.TeamInviteEmailSubject, body)
}
