package security

import "os"

type Settings struct {
	AccessSecret  string
	RefreshSecret string
}

func MustLoadSettings() Settings {
	return Settings{
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
	}
}
