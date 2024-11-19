package logger_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chasehampton/gom/logger"
	"github.com/chasehampton/gom/models"
	"github.com/stretchr/testify/mock"
)

type MockLogInserter struct {
	mock.Mock
}

func (m *MockLogInserter) InsertLogEntry(ctx context.Context, entry models.LogEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func TestLogError(t *testing.T) {
	mockInserter := new(MockLogInserter)
	logger := logger.NewLogger(mockInserter)

	mockEntry := &models.Action{ActionID: 1}
	mockErr := errors.New("Mock error message")

	mockInserter.On("InsertLogEntry", mock.Anything, mock.Anything).Return(nil)

	logger.LogError(mockErr, mockEntry, "")
	mockInserter.AssertExpectations(t)
}

func TestLogMessage(t *testing.T) {
	mockInserter := new(MockLogInserter)
	logger := logger.NewLogger(mockInserter)

	mockEntry := &models.Action{ActionID: 1}
	mockMessage := "Mock message"

	mockInserter.On("InsertLogEntry", mock.Anything, mock.Anything).Return(nil)

	logger.LogMessage(mockMessage, mockEntry, "")
	mockInserter.AssertExpectations(t)
}
