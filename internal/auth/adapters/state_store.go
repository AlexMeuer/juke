package adapters

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/alexmeuer/juke/internal/ports"
)

const (
	stateLength = 64
)

type StateStore struct {
	ports.KeyValueStore[string]
}

func (ss *StateStore) GenerateState(ctx context.Context, ID string) (string, error) {
	buffer := make([]byte, stateLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	state := base64.URLEncoding.EncodeToString(buffer)
	err = ss.Set(ctx, ID, state)
	if err != nil {
		return "", fmt.Errorf("failed to save state: %w", err)
	}
	return state, nil
}

func (ss *StateStore) VerifyState(ctx context.Context, ID, state string) error {
	stored, err := ss.Get(ctx, ID)
	if err != nil {
		return err
	}
	if stored != state {
		return errors.New("state mismatch")
	}
	return nil
}
