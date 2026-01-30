// Package log implements a simple logging package.
package log

import (
	"maps"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"sync"
	"time"
)

var (
	// DefaultFormatter - default log formatter.
	DefaultFormatter = new(DefaultLogFormatter)
	// StdLog - default log
	StdLog = New(os.Stderr, DefaultFormatter, Fields{})
)

// Logger - kira logger.
type Logger struct {
	formatter Formatter
	level     Level
	fields    Fields

	Writer io.Writer
	lock   sync.Mutex
}

const defaultLevel = DebugLevel

// New creates a new Logger.
func New(w io.Writer, f Formatter, fields Fields) *Logger {
	return &Logger{
		Writer:    w,
		formatter: f,
		level:     defaultLevel,
		fields:    fields,
	}
}

// log writes the output for a logging event.
func (l *Logger) log(level Level, msg any) {
	if l.level <= level {
		if err := l.formatter.Format(l, level, msg, time.Now().UTC()); err != nil {
			stdlog.Printf("error logging: %s", err)
		}
	}
}

// SetFormatter sets the Formatter for the logger.
func (l *Logger) SetFormatter(f Formatter) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.formatter = f
}

// SetLevel sets the level severity for the logger.
func (l *Logger) SetLevel(level Level) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.level = level
}

// SetWriter sets the logger Writer
func (l *Logger) SetWriter(w io.Writer) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Writer = w
}

func (l *Logger) WithField(key string, value any) *Logger {
	fields := Fields{}
	maps.Copy(fields, l.fields)
	fields[key] = value

	return New(l.Writer, l.formatter, fields)
}

// Debug calls log to print to the logger.
func (l *Logger) Debug(v ...any) {
	l.log(DebugLevel, fmt.Sprint(v...))
}

// Debugf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(f string, v ...any) {
	l.Debug(fmt.Sprintf(f, v...))
}

// Info calls log to print to the logger.
func (l *Logger) Info(v ...any) {
	l.log(InfoLevel, fmt.Sprint(v...))
}

// Infof calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(f string, v ...any) {
	l.Info(fmt.Sprintf(f, v...))
}

// Warn calls log to print to the logger.
func (l *Logger) Warn(v ...any) {
	l.log(WarnLevel, fmt.Sprint(v...))
}

// Warnf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(f string, v ...any) {
	l.Warn(fmt.Sprintf(f, v...))
}

// Error calls log to print to the logger.
func (l *Logger) Error(v ...any) {
	l.log(ErrorLevel, fmt.Sprint(v...))
}

// Errorf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(f string, v ...any) {
	l.Error(fmt.Sprintf(f, v...))
}

// Fatal calls log to print to the logger.
func (l *Logger) Fatal(v ...any) {
	l.log(FatalLevel, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(f string, v ...any) {
	l.Fatal(fmt.Sprintf(f, v...))
}

// Panic calls log to print to the logger.
func (l *Logger) Panic(v ...any) {
	s := fmt.Sprint(v...)
	l.log(PanicLevel, s)

	panic(s)
}

// Panicf calls l.log to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Panicf(f string, v ...any) {
	l.Panic(fmt.Sprintf(f, v...))
}

// Debug calls log to print to the standard logger.
func Debug(v ...any) {
	StdLog.Debug(v...)
}

// Debugf calls log to print to the standard logger.
func Debugf(f string, v ...any) {
	StdLog.Debug(fmt.Sprintf(f, v...))
}

// Info calls log to print to the standard logger.
func Info(v ...any) {
	StdLog.Info(v...)
}

// Infof calls log to print to the standard logger.
func Infof(f string, v ...any) {
	StdLog.Info(fmt.Sprintf(f, v...))
}

// Warn calls log to print to the standard logger.
func Warn(v ...any) {
	StdLog.Warn(v...)
}

// Warnf calls log to print to the standard logger.
func Warnf(f string, v ...any) {
	StdLog.Warn(fmt.Sprintf(f, v...))
}

// Error calls log to print to the standard logger.
func Error(v ...any) {
	StdLog.Error(v...)
}

// Errorf calls log to print to the standard logger.
func Errorf(f string, v ...any) {
	StdLog.Error(fmt.Sprintf(f, v...))
}

// Fatal calls log to print to the standard logger.
func Fatal(v ...any) {
	StdLog.Fatal(v...)
}

// Fatalf calls log to print to the standard logger.
func Fatalf(f string, v ...any) {
	StdLog.Fatal(fmt.Sprintf(f, v...))
}

// Panic calls log to print to the standard logger.
func Panic(v ...any) {
	StdLog.Panic(v...)
}

// Panicf calls log to print to the standard logger.
func Panicf(f string, v ...any) {
	StdLog.Panic(fmt.Sprintf(f, v...))
}
