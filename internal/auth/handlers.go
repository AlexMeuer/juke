package auth

import (
	"net/http"

	"github.com/alexmeuer/juke/internal/auth/ports"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func NewLoginHandler(cfg *Config, stateGenerator ports.StateGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		log.Info().Str("session ID", session.ID()).Msg("starting login flow")
		verifier := oauth2.GenerateVerifier()

		state, err := stateGenerator.GenerateState(c, session.ID())
		if err != nil {
			log.Err(err).Msg("failed to generate state")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		session.Set("state", state)
		err = session.Save()
		if err != nil {
			log.Err(err).Str("state", state).Msg("failed to save session")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

		log.Info().Str("state", state).Str("verifier", verifier).Str("session ID", session.ID()).Str("URL", url).Msg("redirecting to Spotify")

		c.Redirect(http.StatusFound, url)
	}
}
