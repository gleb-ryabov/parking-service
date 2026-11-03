package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	parkinggrpc "parking/internal/grpc/parking"
)

// App is struct for gRPC server.
type App struct {
	log        *zerolog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New returns app for gRPC server.
func New(
	log *zerolog.Logger,
	parkingService parkinggrpc.Service,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	parkinggrpc.Register(gRPCServer, parkingService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun starts gRPC app or panic.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Fatal().Err(err).Msg("failed run gRPC server")
	}
}

// Run starts gRPC app.
func (a *App) Run() error {
	const op = "app.grpc.Run"

	log := a.log.With().
		Str("op", op).
		Logger()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: listen tcp: %w", op, err)
	}

	log.Info().
		Str("addr", l.Addr().String()).
		Msg("gRPC server is running")

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: gRPC serve: %w", op, err)
	}

	return nil
}

// Stop finishes the gRPC server gracefully.
func (a *App) Stop(ctx context.Context) error {
	const op = "app.grpc.Stop"

	a.log.Info().
		Str("op", op).
		Int("port", a.port).
		Msg("stopping gRPC server")

	done := make(chan struct{})

	go func() {
		a.gRPCServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		a.gRPCServer.Stop()

		return fmt.Errorf("%s: graceful shutdown: %w", op, ctx.Err())
	}
}
