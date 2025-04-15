package log

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

var logger = &Logger{}

type Logger struct {
	sl      *slog.Logger
	lvs     map[string]*slog.LevelVar
	version string
	gitRev  string
}

func Set(v *viper.Viper, opts ...Option) error {
	l, err := New(v, opts...)
	if err != nil {
		return fmt.Errorf("log.Set: %w", err)
	}

	logger = l

	return nil
}

func GetSlog() *slog.Logger {
	return logger.sl
}

func (l *Logger) SetLevel(h, lv string) bool {
	if l.lvs == nil {
		return false
	}

	lvl, ok := l.lvs[h]
	if !ok {
		return false
	}

	lv = strings.ToLower(lv)
	switch lv {
	case "debug":
		lvl.Set(slog.LevelDebug)
	case "info":
		lvl.Set(slog.LevelInfo)
	case "warn":
		lvl.Set(slog.LevelWarn)
	case "error":
		lvl.Set(slog.LevelError)
	default:
		return false
	}

	return true
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

	return &Logger{
		sl:      l.sl.With(args...),
		lvs:     l.lvs,
		version: l.version,
		gitRev:  l.gitRev,
	}
}

func (l *Logger) WithGroup(name string) *Logger {
	if l.sl == nil {
		return l
	}

	return &Logger{
		sl:      l.sl.WithGroup(name),
		lvs:     l.lvs,
		version: l.version,
		gitRev:  l.gitRev,
	}
}

func SetLevel(h, lv string) bool {
	return logger.SetLevel(h, lv)
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
