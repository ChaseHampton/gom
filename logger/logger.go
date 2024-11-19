package logger

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/chasehampton/gom/db/create"
	"github.com/chasehampton/gom/models"
)

type Logger struct {
	inserter create.LogInserter
}

func NewLogger(inserter create.LogInserter) *Logger {
	return &Logger{inserter: inserter}
}

func (i *Logger) LogError(err error, act *models.Action, additional_json string) {
	var stackTrace string
	if err != nil {
		stackTrace = string(debug.Stack())
	}
	entry := models.LogEntry{
		Timestamp:  time.Now(),
		Message:    err.Error(),
		StackTrace: stackTrace,
		ScheduleID: 1,
		ActionID:   act.ActionID,
		Additional: additional_json,
	}
	err = i.inserter.InsertLogEntry(context.TODO(), entry)
	if err != nil {
		log.Printf("Error inserting log entry: %v\n", err)
	}
}

func (i *Logger) LogMessage(message string, act *models.Action, additional_json string) {
	entry := models.LogEntry{
		Timestamp:  time.Now(),
		Message:    message,
		ScheduleID: 1,
		ActionID:   act.ActionID,
		Additional: additional_json,
	}
	err := i.inserter.InsertLogEntry(context.TODO(), entry)
	if err != nil {
		log.Printf("Error inserting log entry: %v\n", err)
	}
}
