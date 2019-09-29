package llog

import (
	"log"
	"testing"
)

func TestNewConsoleLog(t *testing.T) {
	lc := Config{
		OutputFile: "stdout",
		Level:      "debug",
	}

	lg, err := New(lc, log.LstdFlags|log.Lshortfile)
	if err != nil {
		t.Error(err)
	}
	lg.Debug("test %d", 11)
}

func TestNewConsoleLogTimeFormat(t *testing.T) {
	lc := Config{
		OutputFile: "stdout",
		Level:      "debug",
	}

	lg, err := NewWithTimeFormat(lc, log.Lshortfile, "2006-01-02 15:04:05.000")
	if err != nil {
		t.Error(err)
	}
	lg.Debug("test %d", 11)
}
