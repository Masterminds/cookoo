/*
 A log manager and interface for Go.
*/
package log

import (
	"fmt"
	"time"
)

func NewLogManager() *LogManager {
	l := new(LogManager)
	l.Init()
	return l
}

// Describe a logger
type Logger interface {
	// Send the message to the log.
	Send(Entry)
}

// The message passed to the logger.
type Entry struct {
	// The category associated with the message. E.g., debug.
	Category string // Priority of the log message.
	// The message to log.
	Message string
	// The time the message was created
	Created time.Time // Time this message was created.
}

type LogManager struct {
	loggers map[string]Logger
}

func (l *LogManager) Init() {
	l.loggers = make(map[string]Logger)
}

// Get all the loggers currently on the LogManager.
func (l *LogManager) Loggers() map[string]Logger {
	return l.loggers
}

// Retrieve a specific logger by name.
func (l *LogManager) Has(name string) (log Logger, found bool) {
	log, found = l.loggers[name]
	return
}

// Retrieve a specific logger by name.
func (l *LogManager) Logger(name string) Logger {
	return l.loggers[name]
}

// A a logger implementation to write to.
func (l *LogManager) AddLogger(name string, log Logger) *LogManager {
	l.loggers[name] = log
	return l
}

// Remove a logger from use.
func (l *LogManager) RemoveLogger(name string) {
	delete(l.loggers, name)
}

// Log a message.
func (l *LogManager) Log(category string, v ...interface{}) {

	entry := Entry{
		Category: category,
		Message:  fmt.Sprint(v...),
		Created:  time.Now(),
	}

	for _, logger := range l.loggers {
		logger.Send(entry)
	}
}

func (l *LogManager) Logf(category string, format string, v ...interface{}) {
	l.Log(category, fmt.Sprintf(format, v...))
}

func (l *LogManager) Fatal(v ...interface{}) {

}

func (l *LogManager) Fatalf(format string, v ...interface{}) {

}

func (l *LogManager) Fatalln(v ...interface{}) {

}

func (l *LogManager) Flags() int {

}

func (l *LogManager) Output(calldepth int, s string) error {

}

func (l *LogManager) Panic(v ...interface{}) {

}

func (l *LogManager) Panicf(format string, v ...interface{}) {

}

func (l *LogManager) Panicln(v ...interface{}) {

}

func (l *LogManager) Prefix() string {
	return l.Category
}

func (l *LogManager) Print(v ...interface{}) {

}

func (l *LogManager) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Srintf(format, v...))
}

func (l *LogManager) Println(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
}

func (l *LogManager) SetFlags(flag int) {

}

func (l *LogManager) SetPrefix(prefix string) {
	l.Category = prefix
}
