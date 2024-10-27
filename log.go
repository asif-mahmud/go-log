package golog

import (
	"io"
	"log/slog"
	"os"
)

// LogOpt configuration options for settign up built-in loggers
type LogOpt struct {
	writer io.Writer
	json   bool
}

// LogOptFunc definition for configuration builder function
type LogOptFunc func(*LogOpt) *LogOpt

// Setup configures standard library loggers.
// By default, it will write logs in json format and to stderr.
func Setup(funcs ...LogOptFunc) {
	opt := &LogOpt{
		writer: os.Stderr,
		json:   true,
	}

	for _, f := range funcs {
		opt = f(opt)
	}

	var h slog.Handler
	ho := &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}

	if !opt.json {
		h = slog.NewTextHandler(opt.writer, ho)
	} else {
		h = slog.NewJSONHandler(opt.writer, ho)
	}

	l := slog.New(h)
	slog.SetDefault(l)
}

// WithWriter sets log writer to wr.
func WithWriter(wr io.Writer) LogOptFunc {
	return func(lo *LogOpt) *LogOpt {
		lo.writer = wr
		return lo
	}
}

// WithText sets up non-json log message writer.
func WithText() LogOptFunc {
	return func(lo *LogOpt) *LogOpt {
		lo.json = false
		return lo
	}
}
