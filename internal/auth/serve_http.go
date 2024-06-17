package auth

import (
	"fmt"
	"os"
	"strconv"

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

	tokenAesKey, ok := os.LookupEnv("TOKEN_AES_KEY")
	if !ok {
		log.Fatal().Msg("TOKEN_AES_KEY environment variable not set")
	}

	store, err := authAdapters.NewBadgerStore(os.Getenv("BADGER_PATH"), []byte(tokenAesKey))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Badger store")
	}
	defer store.Close()

	spotifyGroup := r.Group("/spotify")
	spotifyGroup.GET("/login", NewLoginHandler(spotifyConfig, store))
	spotifyGroup.GET("/callback", NewCallbackHandler(spotifyConfig, store, store))

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
