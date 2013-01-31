/*
 Provides a console based logger.
*/
package log

// Constructor to create a new ConsoleLogger
func New(cat []string) *ConsoleLogger {
	l := new(ConsoleLogger)
	return l
}

type ConsoleLogger struct{}
