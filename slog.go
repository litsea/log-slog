package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	sentry "github.com/litsea/sentry-slog"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	errEmptyLogConfig    = fmt.Errorf("empty log config")
	errInvalidLogHandler = fmt.Errorf("invalid log handler")
	errNoValidLogHandler = fmt.Errorf("no valid log handler")
	errEmptyFilename     = fmt.Errorf("empty filename")
)

const (
	RFC3339Micro = "2006-01-02T15:04:05.999999Z"
)

const (
	HandlerText   = "text"
	HandlerJSON   = "json"
	HandlerSentry = "sentry"
)

const (
	OutputStdOut = "stdout"
	OutputFile   = "file"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func New(v *viper.Viper, opts ...Option) (*Logger, error) {
	if v == nil {
		return nil, fmt.Errorf("log.New: %w", errEmptyLogConfig)
	}

	l := &Logger{}

	for _, opt := range opts {
		opt(l)
	}

	var (
		cfgHs = v.GetStringSlice("handlers")
		hs    = make([]slog.Handler, 0, len(cfgHs))
		// Log levels for handlers (except sentry)
		lvs = make(map[string]*slog.LevelVar, len(cfgHs))
	)

	var (
		h   slog.Handler
		lv  *slog.LevelVar
		err error
	)

	for _, ch := range cfgHs {
		sub := v.Sub(ch)
		if sub == nil {
			continue
		}

		subH := sub.GetString("handler")
		switch subH {
		case HandlerText:
			h, lv, err = newTextHandler(sub)
		case HandlerJSON:
			h, lv, err = newJSONHandler(sub)
		case HandlerSentry:
			h, err = newSentryHandler(sub, l.gitRev)
			lv = nil
		default:
			err = errInvalidLogHandler
		}

		if err != nil {
			slog.Warn("log.New: parse handler failed", "handler", subH, "err", err)
			continue
		}

		hs = append(hs, h)
		if lv != nil {
			lvs[ch] = lv
		}
	}

	if len(hs) == 0 {
		return nil, fmt.Errorf("log.New: %w", errNoValidLogHandler)
	}

	sl := slog.New(
		slogmulti.Fanout(hs...),
	)
	if l.version != "" {
		sl = sl.With("project.version", l.version)
	}
	l.sl = sl
	l.lvs = lvs

	return l, nil
}

func replaceDateTimeFunc(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "time" && a.Value.Kind() == slog.KindTime {
		a.Value = slog.StringValue(
			a.Value.Time().In(time.UTC).Format(RFC3339Micro),
		)
	}
	return a
}

func newTextHandler(sub *viper.Viper) (slog.Handler, *slog.LevelVar, error) {
	w, err := getWriter(sub)
	if err != nil {
		return nil, nil, fmt.Errorf("log.newTextHandler: %w", err)
	}

	lv := getLevel(sub)

	return slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       lv,
		ReplaceAttr: replaceDateTimeFunc,
	}), lv, nil
}

func newJSONHandler(sub *viper.Viper) (slog.Handler, *slog.LevelVar, error) {
	w, err := getWriter(sub)
	if err != nil {
		return nil, nil, fmt.Errorf("log.newJSONHandler: %w", err)
	}

	lv := getLevel(sub)

	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       lv,
		ReplaceAttr: replaceDateTimeFunc,
	}), lv, nil
}

func getWriter(sub *viper.Viper) (io.Writer, error) {
	o := sub.GetString("output")
	switch o {
	case OutputFile:
		filename := sub.GetString("filename")
		if filename == "" {
			return nil, errEmptyFilename
		}

		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return nil, fmt.Errorf("log.getWriter: %w", err)
		}

		sub.SetDefault("max-mbs", 20)
		sub.SetDefault("max-days", 30)
		sub.SetDefault("max-backups", 10)

		return &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    sub.GetInt("max-mbs"),
			MaxAge:     sub.GetInt("max-days"),
			MaxBackups: sub.GetInt("max-backups"),
			Compress:   true,
		}, nil
	case OutputStdOut:
		return os.Stdout, nil
	default:
		return os.Stdout, nil
	}
}

func getLevel(sub *viper.Viper) *slog.LevelVar {
	lv := new(slog.LevelVar)

	l := strings.ToLower(sub.GetString("level"))
	switch l {
	case LevelDebug:
		lv.Set(slog.LevelDebug)
	case LevelInfo:
		lv.Set(slog.LevelInfo)
	case LevelWarn:
		lv.Set(slog.LevelWarn)
	case LevelError:
		lv.Set(slog.LevelError)
	default:
		lv.Set(slog.LevelInfo)
	}

	return lv
}

func newSentryHandler(sub *viper.Viper, rel string) (slog.Handler, error) {
	h, err := sentry.NewHandler(
		sentry.WithDSN(sub.GetString("dsn")),
		sentry.WithEnvironment(sub.GetString("env")),
		sentry.WithRelease(rel),
		sentry.WithDebug(sub.GetBool("debug")),
		// Other log handlers can set the stacktrace in the log attributes,
		// But for the Sentry handler, we use `sentry.WithAttachStacktrace(true)` to capture the stacktrace,
		// So we need to ignore it here
		// See:
		// github.com/samber/slog-sentry/converter.go#DefaultConverter
		// github.com/samber/slog-common/attributes.go#RemoveEmptyAttrs()
		sentry.WithLogReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "stacktrace" {
				a.Key = ""
			}
			return a
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("log.newSentryHandler: %w", err)
	}

	return h, nil
}
