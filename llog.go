package llog

import (
	"fmt"
	"io"
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
	OutputFile string
	MaxSize    int    // megabytes after which new file is created
	MaxBackups int    // number of backups
	MaxAge     int    // max keep days
	Level      string // log level
}

type Logger struct {
	*log.Logger
	level  int
	Closer io.WriteCloser
}

func (l *Logger) Level() int {
	return l.level
}

func (l *Logger) Close() error {
	if l.Closer != nil {
		return l.Closer.Close()
	}

	return nil
}

func (l *Logger) Finest(format string, args ...interface{}) {
	if l.level > Lfinest {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) Fine(format string, args ...interface{}) {
	if l.level > Lfine {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level > Ldebug {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level > Linfo {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level > Lwarn {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level > Lerr {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) Critical(format string, args ...interface{}) {
	if l.level > Lcritical {
		return
	}

	f, a := callerSource(format, args)
	l.Printf(f, a...)
}

func (l *Logger) SetLevel(level string) error {
	lv, ok := LogLevel[level]
	if !ok {
		return fmt.Errorf("log level is invalid")
	}

	l.level = lv
	l.SetPrefix(levelPrefix[lv])
	return nil
}

func New(lc Config, flag int) (*Logger, error) {
	flag = flag | log.LstdFlags

	lv, ok := LogLevel[lc.Level]
	if !ok {
		return nil, fmt.Errorf("log level is invalid")
	}

	prefix := levelPrefix[lv]

	switch lc.OutputFile {
	case "stdout":
		return &Logger{log.New(os.Stdout, prefix, flag), lv, nil}, nil
	case "stderr":
		return &Logger{log.New(os.Stderr, prefix, flag), lv, nil}, nil
	case "":
		return nil, fmt.Errorf("output file cannot be nil")
	default:
		file, err := os.OpenFile(lc.OutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	lg := &Logger{log.New(nil, prefix, flag), lv, nil}
	lj := &lumberjack.Logger{
		Filename:   lc.OutputFile,
		MaxSize:    lc.MaxSize,    // megabytes after which new file is created
		MaxBackups: lc.MaxBackups, // number of backups
		MaxAge:     lc.MaxAge,     //days
		LocalTime:  true,
	}
	lg.SetOutput(lj)
	lg.Closer = lj

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
