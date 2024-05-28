package adapters_test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockKeyValueStore is a mock implementation of the KeyValueStore interface
type MockKeyValueStore[T any] struct {
	mock.Mock
}

func (m *MockKeyValueStore[T]) Set(ctx context.Context, key string, value T) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockKeyValueStore[T]) Get(ctx context.Context, key string) (T, error) {
	args := m.Called(ctx, key)
	var result T
	if args.Get(0) != nil {
		result = args.Get(0).(T)
	}
	return result, args.Error(1)
}
