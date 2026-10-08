package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// App is struct for HTTP server.
type App struct {
	log        *zerolog.Logger
	app        *gin.Engine
	httpServer *http.Server
}

// New returns app for HTTP server.
func New(
	log *zerolog.Logger,
	port int,
) *App {
	app := gin.New()
	app.Use(gin.Recovery())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app,
	}

	return &App{
		log:        log,
		app:        app,
		httpServer: httpServer,
	}
}

// MustRun starts HTTP app or panic.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Fatal().Err(err).Msg("failed run HTTP server")
	}
}

// Run starts HTTP app.
func (a *App) Run() error {
	const op = "app.http.Run"

	a.log.Info().
		Str("op", op).
		Str("addr", a.httpServer.Addr).
		Msg("HTTP server is running")

	if err := a.httpServer.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			return nil
		default:
			return fmt.Errorf("%s: run HTTP server: %w", op, err)
		}
	}

	return nil
}

// Stop finishes the HTTP server gracefully.
func (a *App) Stop(ctx context.Context) error {
	const op = "app.http.Stop"

	a.log.Info().
		Str("op", op).
		Str("addr", a.httpServer.Addr).
		Msg("stopping HTTP server")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: stop http server: %w", op, err)
	}

	return nil
}
