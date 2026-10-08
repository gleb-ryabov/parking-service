package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"parking/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	appl := app.New(ctx)

	appl.MustRun()

	<-ctx.Done()
	appl.Stop()
}
