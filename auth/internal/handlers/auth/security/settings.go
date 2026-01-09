package security

import (
	"os"
	"strconv"
)

type Settings struct {
	AccessSecret string
	AccessTTL    int
	RefreshTTL   int
}

func MustLoadSettings() Settings {
	attl, err := strconv.Atoi(os.Getenv("SECURITY__ACCESS_LIFETIME_MINUTES"))
	if err != nil {
		panic(err)
	}

	rttl, err := strconv.Atoi(os.Getenv("SECURITY__REFRESH_LIFETIME_DAYS"))
	if err != nil {
		panic(err)
	}

	return Settings{
		AccessSecret: os.Getenv("SECURITY__ACCESS_SECRET"),
		AccessTTL:    attl,
		RefreshTTL:   rttl,
	}
}
