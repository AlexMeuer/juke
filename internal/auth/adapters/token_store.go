package adapters

import (
	"context"

	"github.com/alexmeuer/juke/internal/ports"
	"golang.org/x/oauth2"
)

type TokenStore struct {
	ports.KeyValueStore[*oauth2.Token]
}

func (ts *TokenStore) SaveToken(ctx context.Context, ID string, token *oauth2.Token) error {
	return ts.Set(ctx, ID, token)
}

func (ts *TokenStore) RetrieveToken(ctx context.Context, ID string) (*oauth2.Token, error) {
	return ts.Get(ctx, ID)
}
