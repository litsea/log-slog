package log

import (
	"log/slog"

	"github.com/spf13/viper"
)

var logger *slog.Logger

func Set(v *viper.Viper) error {
	l, err := New(v)
	if err != nil {
		return err
	}

	logger = l

	return nil
}

func Get() *slog.Logger {
	return logger
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func With(args ...any) *slog.Logger {
	return logger.With(args...)
}

func WithGroup(name string) *slog.Logger {
	return logger.WithGroup(name)
}
