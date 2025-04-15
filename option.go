package log

type Option func(l *Logger)

func WithVersion(v string) Option {
	return func(l *Logger) {
		l.version = v
	}
}

func WithGitRev(rev string) Option {
	return func(l *Logger) {
		l.gitRev = rev
	}
}
