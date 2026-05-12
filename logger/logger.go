package logger

import "github.com/gookit/slog"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func DefaultLogger() Logger {
	sl := slog.NewStdLogger(func(sl *slog.SugaredLogger) {
		f := sl.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		sl.ChannelName = "Spoon"
		sl.Level = slog.TraceLevel
	})
	return &log{slog: sl}
}

type log struct {
	slog *slog.SugaredLogger
}

func (l *log) Debug(msg string, args ...any) {
	l.slog.Debugf(msg, args...)
}

func (l *log) Info(msg string, args ...any) {
	l.slog.Infof(msg, args...)
}

func (l *log) Warn(msg string, args ...any) {
	l.slog.Warnf(msg, args...)
}

func (l *log) Error(msg string, args ...any) {
	l.slog.Errorf(msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debugf(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Infof(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warnf(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Errorf(msg, args...)
}
