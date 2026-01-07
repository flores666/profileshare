package mailer

import (
	"fmt"
	"net/smtp"
)

type Service interface {
	Send(to, title, text string) error
}

type service struct {
	settings Settings
	auth     *smtp.Auth
}

func NewMailer(settings Settings) Service {
	return &service{
		settings: settings,
	}
}

func (m *service) Send(to, title, text string) error {
	if m.auth == nil {
		auth := smtp.PlainAuth("", m.settings.Username, m.settings.Password, m.settings.Host)
		m.auth = &auth
	}

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, title, text))

	err := smtp.SendMail(m.settings.Host+":587", *m.auth, m.settings.Username, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}
