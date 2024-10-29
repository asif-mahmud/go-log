package golog

import (
	"io"
	"log/slog"
	"os"
)

// LogOpt configuration options for settign up built-in loggers
type LogOpt struct {
	writer     io.Writer
	json       bool
	handlerOpt *slog.HandlerOptions
	replacers  map[string]AttrReplacerFunc
	attrs      []slog.Attr
}

// LogOptFunc definition for configuration builder function
type LogOptFunc func(*LogOpt) *LogOpt

// Setup configures standard library loggers.
// By default, it will write logs in json format and to stderr.
func Setup(funcs ...LogOptFunc) {
	opt := &LogOpt{
		writer: os.Stderr,
		json:   true,
		handlerOpt: &slog.HandlerOptions{
			AddSource: false,
			Level:     nil,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				return a
			},
		},
		replacers: map[string]AttrReplacerFunc{},
		attrs:     []slog.Attr{},
	}

	for _, f := range funcs {
		opt = f(opt)
	}

	opt.handlerOpt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if replacer, ok := opt.replacers[a.Key]; ok {
			return replacer(a)
		}

		return a
	}

	var h slog.Handler

	if !opt.json {
		h = slog.NewTextHandler(opt.writer, opt.handlerOpt)
	} else {
		h = slog.NewJSONHandler(opt.writer, opt.handlerOpt)
	}

	l := slog.New(h)

	for _, attr := range opt.attrs {
		l = l.With(attr)
	}

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

// WithSource adds caller information in slog's log
func WithSource() LogOptFunc {
	return func(lo *LogOpt) *LogOpt {
		lo.handlerOpt.AddSource = true
		return lo
	}
}

// WithLevel sets log level
func WithLevel(l slog.Level) LogOptFunc {
	return func(lo *LogOpt) *LogOpt {
		lo.handlerOpt.Level = l
		return lo
	}
}

// AttrReplacerFunc function definition for replacing attributes
type AttrReplacerFunc func(slog.Attr) slog.Attr

// WithReplacer adds an attribute replacer to customize formatting
func WithReplacer(key string, replacer AttrReplacerFunc) LogOptFunc {
	return func(lo *LogOpt) *LogOpt {
		lo.replacers[key] = replacer
		return lo
	}
}

// WithAttr adds an attribute in all of the logs
func WithAttr(attr slog.Attr) LogOptFunc {
	return func(lo *LogOpt) *LogOpt {
		lo.attrs = append(lo.attrs, attr)
		return lo
	}
}
