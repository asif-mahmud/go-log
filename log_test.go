package golog_test

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	golog "github.com/asif-mahmud/go-log"
)

func TestDefault(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	golog.Setup(golog.WithWriter(buf))

	log.Println("Hello, World")

	type msg struct {
		Msg string `json:"msg"`
	}
	var m msg

	if err := json.NewDecoder(buf).Decode(&m); err != nil {
		t.Error(err)
	}

	if m.Msg != "Hello, World" {
		t.Errorf("expected: %s, found: %s", "Hello, World", m.Msg)
	}
}
