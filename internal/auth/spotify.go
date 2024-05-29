package auth

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
)

func NewSpotifyConfigFromEnv() (*Config, error) {
	// Don't bail early if we fail to find an env var.
	// Instead, collect all errors and return them all at once.
	var idErr, secretErr, redirectErr error

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		idErr = errors.New("SPOTIFY_CLIENT_ID is required")
	}

	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		secretErr = errors.New("SPOTIFY_CLIENT_SECRET is required")
	}

	redirectURL := os.Getenv("SPOTIFY_REDIRECT_URI")
	if redirectURL == "" {
		redirectErr = errors.New("SPOTIFY_REDIRECT_URI is required")
	}

	return NewConfig(clientID, clientSecret, redirectURL).WithSpotify(),
		errors.Join(idErr, secretErr, redirectErr)
}

func (cfg *Config) WithSpotify() *Config {
	cfg.Scopes = append(cfg.Scopes,
		"user-read-playback-state",
		"user-modify-playback-state",
		"user-read-currently-playing",
	)
	cfg.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://accounts.spotify.com/authorize",
		TokenURL: "https://accounts.spotify.com/api/token",
	}
	return cfg
}
