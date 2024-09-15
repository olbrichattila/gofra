package mail

import (
	"net/smtp"
	"os"
	"strings"
)

func New() Mailer {
	return &Mail{}
}

type Mailer interface {
	Send(string, string, string, string) error
}

type Mail struct{}

func (m *Mail) Send(from, to, subject, message string) error {
	userName := os.Getenv("SMTP_USER_NAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", userName, password, smtpHost)
	composedMessage := m.compose(
		"From: ", from, "\r\n",
		"To: ", to, "\r\n",
		"Subject: ", subject, "\r\n",
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n",
		message,
	)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, composedMessage)
	if err != nil {
		return err
	}

	return nil
}

func (*Mail) compose(parts ...string) []byte {
	cm := strings.Builder{}
	for _, p := range parts {
		cm.WriteString(p)
	}

	return []byte(cm.String())
}
