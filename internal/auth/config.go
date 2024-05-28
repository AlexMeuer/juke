package auth

import (
	"net/http"

	"golang.org/x/oauth2"
)

type Config struct {
	cfg oauth2.Config
}

func NewConfig(clientID, clientSecret, redirectURL string) *Config {
	return &Config{
		cfg: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
		},
	}
}

func (c *Config) Client() *http.Client {
	return c.Client()
}

func (c *Config) AuthorisationURL(state string) string {
	return c.cfg.AuthCodeURL("my-hardcoded-state")
}
