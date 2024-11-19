package db_test

import (
	"errors"
	"os/user"
	"sync"
	"testing"

	"github.com/chasehampton/gom/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUser struct {
	mock.Mock
}

func (m *MockUser) Current() (*user.User, error) {
	args := m.Called()
	return args.Get(0).(*user.User), args.Error(1)
}

func TestGetCurrentUser(t *testing.T) {
	// Reset the once variable to allow multiple tests
	resetOnce := func() {
		db.O = sync.Once{}
	}

	t.Run("success", func(t *testing.T) {
		resetOnce()
		db.CurrentFunc = func() (*user.User, error) {
			return &user.User{Username: "testuser"}, nil
		}
		username := db.GetCurrentUser()
		assert.Equal(t, "testuser", username)
	})

	t.Run("error", func(t *testing.T) {
		resetOnce()
		db.CurrentFunc = func() (*user.User, error) {
			return nil, errors.New("error")
		}
		username := db.GetCurrentUser()
		assert.Equal(t, "unknown", username)
	})
}
