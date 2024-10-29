package golog_test

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"strings"
	"testing"
	"time"

	golog "github.com/asif-mahmud/go-log"
)

func TestDefault(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	golog.Setup(golog.WithWriter(buf))

	type testCase struct {
		msg     string
		logFunc func(string, ...any)
	}

	cases := []testCase{
		{"Hello, World", func(msg string, args ...any) {
			log.Println(msg)
		}},
		{"Hello, World", slog.Info},
	}

	type msg struct {
		Msg string `json:"msg"`
	}

	for _, cs := range cases {
		t.Run(cs.msg, func(t *testing.T) {
			cs.logFunc(cs.msg)

			var m msg

			if err := json.NewDecoder(buf).Decode(&m); err != nil {
				t.Error(err)
			}

			if m.Msg != cs.msg {
				t.Errorf("expected: %s, found: %s", cs.msg, m.Msg)
			}

			buf.Truncate(0)
		})
	}
}

func TestWithSource(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	golog.Setup(golog.WithWriter(buf), golog.WithSource())

	type testCase struct {
		msg     string
		logFunc func(string, ...any)
	}

	cases := []testCase{
		{"Hello, World", slog.Warn},
		{"Hello, World", slog.Info},
	}

	type msg struct {
		Msg    string      `json:"msg"`
		Source slog.Source `json:"source"`
	}

	for _, cs := range cases {
		t.Run(cs.msg, func(t *testing.T) {
			cs.logFunc(cs.msg)

			var m msg

			if err := json.NewDecoder(buf).Decode(&m); err != nil {
				t.Error(err)
			}

			if m.Msg != cs.msg {
				t.Errorf("expected: %s, found: %s", cs.msg, m.Msg)
			}

			if !strings.HasSuffix(m.Source.File, "log_test.go") {
				t.Error("source.file is wrong:", m.Source.File)
			}

			buf.Truncate(0)
		})
	}
}

func TestWithLevel(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	golog.Setup(golog.WithWriter(buf), golog.WithLevel(slog.LevelWarn))

	type testCase struct {
		msg         string
		logFunc     func(string, ...any)
		shouldMatch bool
	}

	cases := []testCase{
		{"Hello, World", slog.Debug, false},
		{"Hello, World", slog.Info, false},
		{"Hello, World", slog.Warn, true},
	}

	type msg struct {
		Msg string `json:"msg"`
	}

	for _, cs := range cases {
		t.Run(cs.msg, func(t *testing.T) {
			cs.logFunc(cs.msg)

			var m msg

			if !cs.shouldMatch && buf.Len() > 0 {
				t.Errorf("should've omitted log, but logged")
			}

			if !cs.shouldMatch {
				return
			}

			if err := json.NewDecoder(buf).Decode(&m); err != nil {
				t.Error(err)
			}

			if cs.shouldMatch && m.Msg != cs.msg {
				t.Errorf("expected: %s, found: %s", cs.msg, m.Msg)
			}

			if !cs.shouldMatch && m.Msg == cs.msg {
				t.Errorf("expected: , found: %s", m.Msg)
			}

			buf.Truncate(0)
		})
	}
}

func TestWithReplacers(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	levels := map[string]string{
		"DEBUG": "debug",
		"INFO":  "info",
	}

	golog.Setup(
		golog.WithWriter(buf),
		golog.WithLevel(slog.LevelDebug),

		// replaces key only
		golog.WithReplacer(slog.TimeKey, func(a slog.Attr) slog.Attr {
			return slog.Attr{
				Key:   "timestamp",
				Value: a.Value,
			}
		}),

		// replaces both key and value
		golog.WithReplacer(slog.LevelKey, func(a slog.Attr) slog.Attr {
			// t.Error(a.Value, a.Value.Kind(), a.Value.Resolve().String())
			if level, ok := levels[a.Value.String()]; ok {
				return slog.Attr{
					Key:   "lvl",
					Value: slog.StringValue(level),
				}
			}

			return slog.Attr{
				Key:   "lvl",
				Value: a.Value,
			}
		}),
	)

	type testCase struct {
		msg     string
		logFunc func(string, ...any)
		level   string
	}

	cases := []testCase{
		{"Hello, World", slog.Debug, "debug"},
		{"Hello, World", slog.Info, "info"},
		{"Hello, World", slog.Warn, "WARN"},
	}

	type msg struct {
		Msg       string    `json:"msg"`
		Timestamp time.Time `json:"timestamp"`
		Level     string    `json:"lvl"`
	}

	for _, cs := range cases {
		t.Run(cs.msg, func(t *testing.T) {
			cs.logFunc(cs.msg)

			var m msg

			if err := json.NewDecoder(buf).Decode(&m); err != nil {
				t.Error(err)
			}

			if cs.msg != m.Msg {
				t.Errorf("expected: %s, found: %s", cs.msg, m.Msg)
			}

			if m.Timestamp.IsZero() {
				t.Errorf("invalid timestamp, got: %v", m.Timestamp)
			}

			if cs.level != m.Level {
				t.Errorf("expected: %s, found: %s", cs.level, m.Level)
			}

			buf.Truncate(0)
		})
	}
}

func TestWithAttribute(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	golog.Setup(golog.WithWriter(buf), golog.WithAttr(slog.Attr{
		Key:   "service",
		Value: slog.StringValue("my-service"),
	}))

	type testCase struct {
		msg     string
		logFunc func(string, ...any)
		service string
	}

	cases := []testCase{
		{"Hello, World", slog.Warn, "my-service"},
		{"Hello, World", slog.Info, "my-service"},
	}

	type msg struct {
		Msg     string `json:"msg"`
		Service string `json:"service"`
	}

	for _, cs := range cases {
		t.Run(cs.msg, func(t *testing.T) {
			cs.logFunc(cs.msg)

			var m msg

			if err := json.NewDecoder(buf).Decode(&m); err != nil {
				t.Error(err)
			}

			if m.Msg != cs.msg {
				t.Errorf("expected: %s, found: %s", cs.msg, m.Msg)
			}

			if m.Service != cs.service {
				t.Errorf("expected: %s, found: %s", cs.service, m.Service)
			}

			buf.Truncate(0)
		})
	}
}
