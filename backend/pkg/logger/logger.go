package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger(levelStr string) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{logger: l}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.logger.Error().Fields(fields).Msg(msg)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatal().Fields(fields).Msg(msg)
	os.Exit(1)
}
