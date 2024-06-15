package auth

import (
	"errors"
	"net/http"

	"github.com/alexmeuer/juke/internal/auth/ports"
	"github.com/alexmeuer/juke/pkg/spotify"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"github.com/segmentio/ksuid"
)

// TODO: Create a wrapper type of session that includes helpers for
// getting/setting the login flow ID and verifier, etc

func NewLoginHandler(cfg *Config, stateGenerator ports.StateGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		loginFlowID := ksuid.New().String()
		session.Set("Login Flow ID", loginFlowID)

		log.Info().Str("login flow ID", loginFlowID).Msg("starting login flow")
		verifier := oauth2.GenerateVerifier()

		state, err := stateGenerator.GenerateState(c, loginFlowID)
		if err != nil {
			log.Err(err).Str("login flow ID", loginFlowID).Msg("failed to generate state")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		session.Set("verifier", verifier)
		err = session.Save()
		if err != nil {
			log.Err(err).Str("login flow ID", loginFlowID).Msg("failed to save session")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

		log.Info().
			Str("state", state).
			Str("verifier", verifier).
			Str("login flow ID", loginFlowID).
			Str("URL", url).Msg("redirecting to Spotify")

		c.Redirect(http.StatusFound, url)
	}
}

func NewCallbackHandler(cfg *Config, stateVerifier ports.StateVerifier, saver ports.TokenSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		loginFlowID, ok := session.Get("Login Flow ID").(string)
		if !ok {
			log.Error().Msg("login flow ID not found in session")
			c.AbortWithError(http.StatusBadRequest, errors.New("login flow ID not found in session"))
			return
		}

		verifier, ok := session.Get("verifier").(string)
		if !ok {
			log.Error().Str("login flow ID", loginFlowID).Msg("verifier not found in session")
			c.AbortWithError(http.StatusBadRequest, errors.New("verifier not found in session"))
			return
		}

		// Ensure the states match before proceeding.
		state := c.Query("state")
		if err := stateVerifier.VerifyState(c, loginFlowID, state); err != nil {
			log.Err(err).
				Str("state", state).
				Str("login flow ID", loginFlowID).
				Msg("failed to verify state")
			c.AbortWithError(http.StatusConflict, err)
			return
		}

		// Exchange the code for a token.
		code := c.Query("code")
		token, err := cfg.Exchange(c, code, oauth2.VerifierOption(verifier))
		if err != nil {
			log.Err(err).
				Str("code", code).
				Str("login flow ID", loginFlowID).
				Msg("failed to exchange code for token")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Use the token to fetch the user's information from Spotify.
		// FIXME: This is the only part of the oauth flow here that is specific to Spotify.
		// We should find a way to do this for any config that we're given!
		client := spotify.New(cfg.Client(c, token))
		me, err := client.Me(c)
		if err != nil {
			log.Err(err).
				Str("login flow ID", loginFlowID).
				Msg("failed to fetch user information")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Info().Interface("me", me).Msg("me")

		// Save the token.
		if err := saver.SaveToken(c, me.ID, token); err != nil {
			log.Err(err).
				Str("login flow ID", loginFlowID).
				Str("user ID", me.ID).
				Msg("failed to save token")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
