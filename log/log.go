package log

import (
	"fmt"
	"log"
)

type Log struct {
	name string
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func NewLogger(name string) Log {
	return Log{name: name}
}

func (l Log) nameFormat() string {
	return fmt.Sprintf("[%s]", l.name)
}

func (l Log) Logf(formatIn string, args ...interface{}) {
	format := fmt.Sprintf("%s %s", l.nameFormat(), formatIn)
	log.Printf(format, args...)
}

func (l Log) Logln(args ...interface{}) {
	l.Logf("%s", fmt.Sprint(args...))
}

func (l Log) Println(args ...interface{}) {
	l.Logln(args...)
}

func (l Log) Printf(format string, args ...interface{}) {
	l.Logf(format, args...)
}
