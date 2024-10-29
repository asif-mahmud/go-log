package golog_test

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"testing"

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
