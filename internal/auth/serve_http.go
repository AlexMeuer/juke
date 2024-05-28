package auth

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ServeHTTP() {
	r := gin.Default()
	r.Use(ginzerolog.Logger("gin"))

	log.Fatal().Msg("not implemented")
	// spotifyConfig, err := NewSpotifyConfigFromEnv()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to create Spotify config")
	// }

	// spotifyGroup := r.Group("/spotify")
	// spotifyGroup.GET("/login", func (c *gin.Context) {
	// 	c.Redirect(http.StatusFound,
	// })
	// spotifyGroup.GET("/callback", SpotifyCallback)
}
