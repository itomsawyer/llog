package llog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/natefinch/lumberjack"
)

var LogLevel = map[string]int{
	"finest":   0,
	"fine":     1,
	"debug":    2,
	"trace":    3,
	"info":     4,
	"warn":     5,
	"error":    6,
	"critical": 7,
}

const (
	Lfinest = iota
	Lfine
	Ldebug
	Ltrace
	Linfo
	Lwarn
	Lerr
	Lcritical
)

var levelPrefix = [Lcritical + 1]string{"[S] ", "[F] ", "[D] ", "[T] ", "[I] ", "[W] ", "[E] ", "[C] "}

type Config struct {
	OutputFile string `toml:"file"`
	MaxSize    int    `toml:"max_size_mb"` // megabytes after which new file is created
	MaxBackups int    `toml:"max_backups"` // number of backups
	MaxAge     int    `toml:"max_age"`     // max keep days
	Level      string `toml:"level"`       // log level
}

func DefaultConfig() Config {
	return Config{OutputFile: "stdout", Level: "debug"}
}

type Logger struct {
	*log.Logger
	level int
	io.WriteCloser
}

func (l *Logger) Writer() io.Writer {
	return l.WriteCloser
}

func (l *Logger) Level() int {
	return l.level
}

func (l *Logger) Close() error {
	if l.WriteCloser != nil {
		return l.WriteCloser.Close()
	}

	return nil
}

func (l *Logger) Finest(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Lfinest {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Lfinest]+format, args...))
}

func (l *Logger) Fine(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Lfine {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Lfine]+format, args...))
}

func (l *Logger) Trace(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Ltrace {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Ltrace]+format, args...))
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Ldebug {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Ldebug]+format, args...))
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Linfo {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Linfo]+format, args...))
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Lwarn {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Lerr]+format, args...))
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Lerr {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Lerr]+format, args...))
}

func (l *Logger) Critical(format string, args ...interface{}) {
	if l.Logger == nil || l.level > Lcritical {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[Lcritical]+format, args...))
}

func (l *Logger) SetLevel(level string) error {
	lv, ok := LogLevel[level]
	if !ok {
		return fmt.Errorf("log level is invalid")
	}

	l.level = lv
	return nil
}

func (l *Logger) Logf(level int, format string, args ...interface{}) {
	if l.Logger == nil || l.level > level {
		return
	}

	l.Output(2, fmt.Sprintf(levelPrefix[level]+format, args...))
}

func NewDefaultLogger() *Logger {
	l, _ := New(Config{OutputFile: "stdout", Level: "debug"}, log.Lshortfile)
	return l
}

func NewEmptyLogger() *Logger {
	l, _ := New(Config{OutputFile: "nil", Level: "debug"}, log.Lshortfile)
	return l
}

func New(lc Config, flag int) (*Logger, error) {
	if flag == 0 {
		flag = log.LstdFlags
	}

	lv, ok := LogLevel[lc.Level]
	if !ok {
		return nil, fmt.Errorf("log level is invalid")
	}

	switch lc.OutputFile {
	case "stdout":
		return &Logger{log.New(os.Stdout, "", flag), lv, nil}, nil
	case "stderr":
		return &Logger{log.New(os.Stderr, "", flag), lv, nil}, nil
	case "nil":
		return &Logger{log.New(ioutil.Discard, "", flag), lv, nil}, nil
	case "":
		return nil, fmt.Errorf("output file cannot be nil")
	default:
		file, err := os.OpenFile(lc.OutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	lg := &Logger{log.New(nil, "", flag), lv, nil}
	lj := &lumberjack.Logger{
		Filename:   lc.OutputFile,
		MaxSize:    lc.MaxSize,    // megabytes after which new file is created
		MaxBackups: lc.MaxBackups, // number of backups
		MaxAge:     lc.MaxAge,     //days
		LocalTime:  true,
	}
	lg.SetOutput(lj)
	lg.WriteCloser = lj

	return lg, nil
}

func callerSource(format string, args []interface{}) (string, []interface{}) {
	// Determine caller func
	_, file, lineno, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		lineno = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	src := fmt.Sprintf("[%s:%d]", file, lineno)
	format = "%s " + format
	a := []interface{}{src}
	a = append(a, args...)
	return format, a
}
