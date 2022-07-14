package log

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

type Log struct {
	name string
	log.Logger
}

var relativePrefix string

func init() {
	_, file, _, _ := runtime.Caller(0)
	relativePrefix = filepath.Dir(file)
	relativePrefix = relativePrefix[:len(relativePrefix)-4]
}

func NewLogger(name string) Log {
	return Log{name: name}
}

func (l *Log) nameFormat() string {
	_, file, line, _ := runtime.Caller(2)
	if strings.HasPrefix(filepath.Dir(file), relativePrefix) {
		file = file[len(relativePrefix)+1:]
	}
	return fmt.Sprintf("%s:%d, [%s]", file, line, l.name)
}

func (l *Log) Logf(formatIn string, args ...interface{}) {
	format := fmt.Sprintf("%s %s", l.nameFormat(), formatIn)
	log.Printf(format, args...)
}

func (l *Log) Logln(args ...interface{}) {
	log.Printf("%s %s", l.nameFormat(), fmt.Sprint(args...))
}

func (l *Log) Printf(formatIn string, args ...interface{}) {
	format := fmt.Sprintf("%s %s", l.nameFormat(), formatIn)
	log.Printf(format, args...)
}

func (l *Log) Println(args ...interface{}) {
	log.Printf("%s %s", l.nameFormat(), fmt.Sprint(args...))
}
