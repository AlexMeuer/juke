package adapters_test

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/alexmeuer/juke/internal/auth/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateState(t *testing.T) {
	ctx := context.Background()
	ID := "test-id"

	t.Run("Success", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[string])
		ss := &adapters.StateStore{
			KeyValueStore: mockStore,
		}

		mockStore.On("Set", ctx, ID, mock.Anything).Return(nil)

		state, err := ss.GenerateState(ctx, ID)

		t.Run("NoError", func(t *testing.T) {
			assert.NoError(t, err)
		})

		t.Run("NotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, state)
		})

		t.Run("CorrectLength", func(t *testing.T) {
			assert.Equal(t, base64.URLEncoding.EncodedLen(64), len(state))
		})

		t.Run("SetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Set", ctx, ID, state)
		})
	})

	t.Run("SetError", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[string])
		ss := &adapters.StateStore{
			KeyValueStore: mockStore,
		}

		expectedError := errors.New("set error")
		mockStore.On("Set", ctx, ID, mock.Anything).Return(expectedError)

		state, err := ss.GenerateState(ctx, ID)

		t.Run("ErrorOccurred", func(t *testing.T) {
			assert.Error(t, err)
		})

		t.Run("StateIsEmpty", func(t *testing.T) {
			assert.Equal(t, "", state)
		})

		t.Run("ErrorMessageContains", func(t *testing.T) {
			assert.Contains(t, err.Error(), "failed to save state")
		})
	})
}

func TestVerifyState(t *testing.T) {
	ctx := context.Background()
	ID := "test-id"
	state := "test-state"

	t.Run("Success", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[string])
		ss := &adapters.StateStore{
			KeyValueStore: mockStore,
		}

		mockStore.On("Get", ctx, ID).Return(state, nil)

		err := ss.VerifyState(ctx, ID, state)

		t.Run("NoError", func(t *testing.T) {
			assert.NoError(t, err)
		})

		t.Run("GetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Get", ctx, ID)
		})
	})

	t.Run("GetError", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[string])
		ss := &adapters.StateStore{
			KeyValueStore: mockStore,
		}

		expectedError := errors.New("get error")
		mockStore.On("Get", ctx, ID).Return("", expectedError)

		err := ss.VerifyState(ctx, ID, "")

		t.Run("ErrorOccurred", func(t *testing.T) {
			assert.Error(t, err)
		})

		t.Run("ErrorMatches", func(t *testing.T) {
			assert.Equal(t, expectedError, err)
		})

		t.Run("GetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Get", ctx, ID)
		})
	})

	t.Run("StateMismatch", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[string])
		ss := &adapters.StateStore{
			KeyValueStore: mockStore,
		}

		storedState := "stored-state"
		givenState := "given-state"

		mockStore.On("Get", ctx, ID).Return(storedState, nil)

		err := ss.VerifyState(ctx, ID, givenState)

		t.Run("ErrorOccurred", func(t *testing.T) {
			assert.Error(t, err)
		})

		t.Run("ErrorMessageMatches", func(t *testing.T) {
			assert.Equal(t, "state mismatch", err.Error())
		})

		t.Run("GetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Get", ctx, ID)
		})
	})
}
