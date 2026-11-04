package main

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"parking/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log.Info().Msg("configuration loaded")

	m, err := migrate.New("file://migrations", cfg.DbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize migrations")
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("no migrations to apply")

			return
		}

		log.Fatal().Err(err).Msg("failed to apply migrations")
	}

	log.Info().Msg("migrations applied")
}
