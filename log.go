package log

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

var logger = &Logger{}

type Logger struct {
	sl *slog.Logger
}

func Set(v *viper.Viper) error {
	l, err := New(v)
	if err != nil {
		return fmt.Errorf("log.Set: %w", err)
	}

	logger = &Logger{sl: l}

	return nil
}

func GetSlog() *slog.Logger {
	return logger.sl
}

func (l *Logger) Debug(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.sl.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.sl.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.sl.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.sl.Error(msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	if logger.sl == nil {
		return logger
	}

	l.sl = l.sl.With(args...)
	return l
}

func (l *Logger) WithGroup(name string) *Logger {
	if l.sl == nil {
		return l
	}

	l.sl = logger.sl.WithGroup(name)
	return l
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

func With(args ...any) *Logger {
	return logger.With(args...)
}

func WithGroup(name string) *Logger {
	return logger.WithGroup(name)
}
