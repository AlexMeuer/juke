package adapters_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexmeuer/juke/internal/auth/adapters"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestSaveToken(t *testing.T) {
	ctx := context.Background()
	ID := "test-id"
	token := &oauth2.Token{AccessToken: "test-token"}

	t.Run("Success", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[*oauth2.Token])
		ts := &adapters.TokenStore{
			KeyValueStore: mockStore,
		}

		mockStore.On("Set", ctx, ID, token).Return(nil)

		err := ts.SaveToken(ctx, ID, token)

		t.Run("NoError", func(t *testing.T) {
			assert.NoError(t, err)
		})

		t.Run("SetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Set", ctx, ID, token)
		})
	})

	t.Run("SetError", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[*oauth2.Token])
		ts := &adapters.TokenStore{
			KeyValueStore: mockStore,
		}

		expectedError := errors.New("set error")
		mockStore.On("Set", ctx, ID, token).Return(expectedError)

		err := ts.SaveToken(ctx, ID, token)

		t.Run("ErrorOccurred", func(t *testing.T) {
			assert.Error(t, err)
		})

		t.Run("ErrorMatches", func(t *testing.T) {
			assert.Equal(t, expectedError, err)
		})
	})
}

func TestRetrieveToken(t *testing.T) {
	ctx := context.Background()
	ID := "test-id"
	token := &oauth2.Token{AccessToken: "test-token"}

	t.Run("Success", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[*oauth2.Token])
		ts := &adapters.TokenStore{
			KeyValueStore: mockStore,
		}

		mockStore.On("Get", ctx, ID).Return(token, nil)

		returnedToken, err := ts.RetrieveToken(ctx, ID)

		t.Run("NoError", func(t *testing.T) {
			assert.NoError(t, err)
		})

		t.Run("TokenMatches", func(t *testing.T) {
			assert.Equal(t, token, returnedToken)
		})

		t.Run("GetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Get", ctx, ID)
		})
	})

	t.Run("GetError", func(t *testing.T) {
		mockStore := new(MockKeyValueStore[*oauth2.Token])
		ts := &adapters.TokenStore{
			KeyValueStore: mockStore,
		}

		expectedError := errors.New("get error")
		mockStore.On("Get", ctx, ID).Return(nil, expectedError)

		returnedToken, err := ts.RetrieveToken(ctx, ID)

		t.Run("ErrorOccurred", func(t *testing.T) {
			assert.Error(t, err)
		})

		t.Run("ErrorMatches", func(t *testing.T) {
			assert.Equal(t, expectedError, err)
		})

		t.Run("TokenIsNil", func(t *testing.T) {
			assert.Nil(t, returnedToken)
		})

		t.Run("GetCalled", func(t *testing.T) {
			mockStore.AssertCalled(t, "Get", ctx, ID)
		})
	})
}
