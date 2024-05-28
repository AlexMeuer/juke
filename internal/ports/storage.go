package ports

import "context"

type KeyValueStore[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T) error
}
