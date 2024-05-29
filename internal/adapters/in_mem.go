package adapters

import (
	"context"
	"errors"
)

type InMemoryKeyValueStore[T any] struct {
	store map[string]T
}

func NewInMemoryKeyValueStore[T any]() *InMemoryKeyValueStore[T] {
	return &InMemoryKeyValueStore[T]{
		store: make(map[string]T),
	}
}

func (kvs *InMemoryKeyValueStore[T]) Get(ctx context.Context, key string) (T, error) {
	value, ok := kvs.store[key]
	if !ok {
		return *new(T), errors.New("key not found")
	}
	return value, nil
}

func (kvs *InMemoryKeyValueStore[T]) Set(ctx context.Context, key string, value T) error {
	kvs.store[key] = value
	return nil
}
