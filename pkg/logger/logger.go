package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger setup and returns logger.
func NewLogger(lvl string, filename string, age, backups int) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "_t"
	zerolog.LevelFieldName = "_l"
	zerolog.CallerFieldName = "_c"
	zerolog.MessageFieldName = "_m"
	zerolog.ErrorFieldName = "_e"
	zerolog.ErrorStackFieldName = "_s"

	levelWithFormat := strings.ToLower(
		strings.TrimSpace(lvl),
	)

	level, err := zerolog.ParseLevel(levelWithFormat)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse logging level")
	}

	zerolog.SetGlobalLevel(level)

	l := log.With().Logger()
	if filename != "" {
		l = l.With().Str("file", filename).Logger()
	}

	logger := zerolog.New(setWriter(filename, age, backups)).With().Timestamp().Caller().Stack().Logger()

	l.Info().Msg("logger set")

	return &logger
}

// setWriter setup writer for logger.
func setWriter(filename string, age, backups int) io.Writer {
	if filename == "" {
		return zerolog.ConsoleWriter{Out: os.Stdout}
	}

	return &lumberjack.Logger{
		Filename:   filename,
		MaxAge:     age,
		MaxBackups: backups,
		Compress:   true,
	}
}
