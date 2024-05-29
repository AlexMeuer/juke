package auth

import (
	"golang.org/x/oauth2"
)

type Config struct {
	oauth2.Config
}

func NewConfig(clientID, clientSecret, redirectURL string) *Config {
	return &Config{
		oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
		},
	}
}
