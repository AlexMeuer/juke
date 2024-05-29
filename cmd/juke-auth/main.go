package main

import (
	"os"

	"github.com/alexmeuer/juke/internal/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := auth.ServeHTTP(); err != nil {
		log.Fatal().Err(err).Msg("failed to serve HTTP")
	}
}
