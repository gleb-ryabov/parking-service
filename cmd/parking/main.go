package main

import (
	"parking/internal/config"
	"parking/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(
		cfg.Log.Level,
		cfg.Log.FileName,
		cfg.Log.MaxAge,
		cfg.Log.MaxBackups,
	)
	_ = log
}
