package mailer

import "os"

type Settings struct {
	Host     string
	Username string
	Password string
}

func MustLoad() Settings {
	return Settings{
		Host:     os.Getenv("SMTP__HOST"),
		Username: os.Getenv("SMTP__USERNAME"),
		Password: os.Getenv("SMTP__PASSWORD"),
	}
}
