package app

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	grpcapp "parking/internal/app/grpc"
	httpapp "parking/internal/app/http"
	"parking/internal/config"
	"parking/internal/service/parking"
	"parking/pkg/logger"
)

const (
	defaultTimeoutStopApp = time.Second * 10
)

// App is struct for instant application.
type App struct {
	log     *zerolog.Logger
	httpApp *httpapp.App
	grpcApp *grpcapp.App
}

// New creates instant application.
func New(ctx context.Context) *App {
	cfg := config.MustLoad()
	log.Info().Msg("configuration loaded")

	l := logger.NewLogger(
		cfg.Log.Level,
		cfg.Log.FileName,
		cfg.Log.MaxAge,
		cfg.Log.MaxBackups,
	)
	l.Info().Str("level", l.GetLevel().String()).Msg("logger initialized")

	// TODO : add repo
	parkingS := parking.New(l, nil)

	grpcApp := grpcapp.New(l, parkingS, cfg.GRPC.Port)
	httpApp := httpapp.New(l, cfg.HTTP.Port)

	return &App{
		log:     l,
		httpApp: httpApp,
		grpcApp: grpcApp,
	}
}

// MustRun starts application server or panic.
func (a *App) MustRun() {
	go a.grpcApp.MustRun()
	go a.httpApp.MustRun()

	a.log.Info().Msg("app is running")
}

// Stop finishes application gracefully.
func (a *App) Stop() {
	const op = "app.Stop"

	a.log.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutStopApp)
	defer cancel()

	if err := a.httpApp.Stop(ctx); err != nil {
		a.log.Error().Err(err).Msg("failed shutdown http server")
	}

	if err := a.grpcApp.Stop(ctx); err != nil {
		a.log.Error().Err(err).Msg("failed shutdown grpc server")
	}

	a.log.Info().Msg("app stopped")
}
