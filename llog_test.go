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
