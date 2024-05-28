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

// MockKeyValueStore is a mock implementation of the KeyValueStore interface
type MockKeyValueStore struct {
	mock.Mock
}

func (m *MockKeyValueStore) Set(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockKeyValueStore) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func TestGenerateState(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	ss := &adapters.StateStore{
		KeyValueStore: mockStore,
	}

	ctx := context.Background()
	ID := "test-id"

	mockStore.On("Set", ctx, ID, mock.Anything).Return(nil)

	state, err := ss.GenerateState(ctx, ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, state)
	assert.Equal(t, base64.URLEncoding.EncodedLen(64), len(state))

	mockStore.AssertCalled(t, "Set", ctx, ID, state)
}

func TestGenerateState_SetError(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	ss := &adapters.StateStore{
		KeyValueStore: mockStore,
	}

	ctx := context.Background()
	ID := "test-id"

	expectedError := errors.New("set error")
	mockStore.On("Set", ctx, ID, mock.Anything).Return(expectedError)

	state, err := ss.GenerateState(ctx, ID)

	assert.Error(t, err)
	assert.Equal(t, "", state)
	assert.Contains(t, err.Error(), "failed to save state")
}

func TestVerifyState(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	ss := &adapters.StateStore{
		KeyValueStore: mockStore,
	}

	ctx := context.Background()
	ID := "test-id"
	state := "test-state"

	mockStore.On("Get", ctx, ID).Return(state, nil)

	err := ss.VerifyState(ctx, ID, state)

	assert.NoError(t, err)
	mockStore.AssertCalled(t, "Get", ctx, ID)
}

func TestVerifyState_GetError(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	ss := &adapters.StateStore{
		KeyValueStore: mockStore,
	}

	ctx := context.Background()
	ID := "test-id"

	expectedError := errors.New("get error")
	mockStore.On("Get", ctx, ID).Return("", expectedError)

	err := ss.VerifyState(ctx, ID, "")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockStore.AssertCalled(t, "Get", ctx, ID)
}

func TestVerifyState_StateMismatch(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	ss := &adapters.StateStore{
		KeyValueStore: mockStore,
	}

	ctx := context.Background()
	ID := "test-id"
	storedState := "stored-state"
	givenState := "given-state"

	mockStore.On("Get", ctx, ID).Return(storedState, nil)

	err := ss.VerifyState(ctx, ID, givenState)

	assert.Error(t, err)
	assert.Equal(t, "state mismatch", err.Error())
	mockStore.AssertCalled(t, "Get", ctx, ID)
}
