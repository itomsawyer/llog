package llog

import (
	"testing"
)

func TestNewConsoleLog(t *testing.T) {
	lc := Config{
		OutputFile: "stdout",
		Level:      "debug",
	}

	lg, err := New(lc, 0)
	if err != nil {
		t.Error(err)
	}
	lg.Debug("test %d", 11)
}
