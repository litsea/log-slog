package log

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var logger = &Logger{}

type Logger struct {
	sl        *slog.Logger
	lvs       map[string]*slog.LevelVar
	version   string
	gitRev    string
	skip      int
	addSource bool
}

func Set(v *viper.Viper, opts ...Option) error {
	l, err := New(v, opts...)
	if err != nil {
		return fmt.Errorf("log.Set: %w", err)
	}

	logger = l

	return nil
}

func Get() *Logger {
	return logger
}

func GetSlog() *slog.Logger {
	return logger.sl
}

func (l *Logger) GetSlog() *slog.Logger {
	return l.sl
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

func (l *Logger) log(ctx context.Context, lv slog.Level, msg string, args ...any) {
	if !l.sl.Enabled(ctx, lv) {
		return
	}

	var pc uintptr
	if l.addSource {
		// skip [runtime.Callers, this function, this function's caller]
		skip := 3
		if l.skip > 0 {
			skip = l.skip
		}
		var pcs [1]uintptr
		runtime.Callers(skip, pcs[:])
		pc = pcs[0]
	}

	r := slog.NewRecord(time.Now(), lv, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.sl.Handler().Handle(ctx, r)
}

func (l *Logger) Debug(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.log(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.log(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.log(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	if l.sl == nil {
		return
	}

	l.log(context.Background(), slog.LevelError, msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	if logger.sl == nil {
		return logger
	}

	return &Logger{
		sl:        l.sl.With(args...),
		lvs:       l.lvs,
		version:   l.version,
		gitRev:    l.gitRev,
		skip:      l.skip,
		addSource: l.addSource,
	}
}

func (l *Logger) WithGroup(name string) *Logger {
	if l.sl == nil {
		return l
	}

	return &Logger{
		sl:        l.sl.WithGroup(name),
		lvs:       l.lvs,
		version:   l.version,
		gitRev:    l.gitRev,
		skip:      l.skip,
		addSource: l.addSource,
	}
}

func SetLevel(h, lv string) bool {
	return logger.SetLevel(h, lv)
}

func Debug(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelError, msg, args...)
}

func With(args ...any) *Logger {
	return logger.With(args...)
}

func WithGroup(name string) *Logger {
	return logger.WithGroup(name)
}
