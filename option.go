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

func WithSkipLevel(skip int) Option {
	return func(l *Logger) {
		if skip > 0 {
			l.skip = skip
		}
	}
}

func WithAddSource(v bool) Option {
	return func(l *Logger) {
		l.addSource = v
	}
}
