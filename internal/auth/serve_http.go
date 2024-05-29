package auth

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alexmeuer/juke/internal/adapters"
	authAdapters "github.com/alexmeuer/juke/internal/auth/adapters"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ServeHTTP() error {
	r := gin.Default()
	r.Use(ginzerolog.Logger("gin"))

	sessionStore := cookie.NewStore(sessionSecret())
	r.Use(sessions.Sessions("juke", sessionStore))

	spotifyConfig, err := NewSpotifyConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Spotify config")
	}

	spotifyGroup := r.Group("/spotify")
	spotifyGroup.GET("/login", NewLoginHandler(spotifyConfig, &authAdapters.StateStore{
		KeyValueStore: adapters.NewInMemoryKeyValueStore[string](),
	}))
	// spotifyGroup.GET("/callback", SpotifyCallback)

	return r.Run(fmt.Sprintf(":%d", port()))
}

// TODO: Use pairs of auth- and excrypt- keys.
//
// See cookie.NewStore() for more information.
func sessionSecret() []byte {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal().Msg("SESSION_SECRET environment variable not set")
	}

	return []byte(secret)
}

func port() uint16 {
	port := os.Getenv("PORT")
	if port == "" {
		log.Warn().Msg("PORT environment variable not set, defaulting to 8080")
		return 8080
	}

	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse PORT environment variable")
	}

	return uint16(portInt)
}
