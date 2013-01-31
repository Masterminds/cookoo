package log

import (
	"testing"
)

type TestLogger struct{}

func (l *TestLogger) Send(entry Entry) {

}

func TestAdding(t *testing.T) {
	manager := NewLogManager()
	logger := new(TestLogger)

	manager.AddLogger("foo", logger)

	logger2, _ := manager.Has("foo")

	if logger != logger2 {
		t.Error("The logger retrieved from Has is incorrect.")
	}

	logger3 := new(TestLogger)

	manager.AddLogger("bar", logger3)

	loggers := manager.Loggers()

	if loggers["foo"] != logger || loggers["bar"] != logger3 {
		t.Error("Expected logger to be present.")
	}

	logger4 := manager.Logger("bar")

	if logger3 != logger4 {
		t.Error("The logger retrieved from Logger is incorrect.")
	}
}

func TestRemoving(t *testing.T) {
	manager := NewLogManager()
	logger := new(TestLogger)

	manager.AddLogger("foo", logger)

	manager.RemoveLogger("foo")

	_, found := manager.Has("foo")

	if found {
		t.Error("Logger found but expected not to be there.")
		// The Ice Age movies are great background distrations while coding.
	}
}
