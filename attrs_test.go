package golog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	golog "github.com/asif-mahmud/go-log"
	"github.com/stretchr/testify/assert"
)

func TestAttrs(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	golog.Setup(golog.WithWriter(buf), golog.WithLevel(slog.LevelDebug))

	slog.Debug(
		"Hello, World",
		golog.Path("/"),
		golog.Query(url.Values{"q": {"search"}}),
		golog.Method("GET"),
		golog.Extra(map[string]string{"data": "dummy"}),
		golog.UserAgent("firefox"),
		golog.Ip("::1234"),
		golog.Status(http.StatusOK),
		golog.Latency(1*time.Second),
		golog.Length(100),
	)

	type msg struct {
		Msg       string            `json:"msg"`
		Path      string            `json:"path"`
		Query     url.Values        `json:"query"`
		Method    string            `json:"method"`
		Extra     map[string]string `json:"extra"`
		UserAgent string            `json:"useragent"`
		Ip        string            `json:"ip"`
		Status    int               `json:"status"`
		Latency   float64           `json:"latency"`
		Length    int               `json:"length"`
	}

	expected := msg{
		Msg:       "Hello, World",
		Path:      "/",
		Query:     map[string][]string{"q": {"search"}},
		Method:    "GET",
		Extra:     map[string]string{"data": "dummy"},
		UserAgent: "firefox",
		Ip:        "::1234",
		Status:    http.StatusOK,
		Latency:   1.0,
		Length:    100,
	}

	var found msg
	if err := json.NewDecoder(buf).Decode(&found); err != nil {
		t.Error(err)
	}

	assert := assert.New(t)

	assert.Equal(expected, found)
}
