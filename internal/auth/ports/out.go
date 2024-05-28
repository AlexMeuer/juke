package ports

import (
	"context"

	"golang.org/x/oauth2"
)

type TokenSaver interface {
	SaveToken(ctx context.Context, ID string, token *oauth2.Token) error
}

type TokenRetriever interface {
	RetrieveToken(ctx context.Context, ID string) (*oauth2.Token, error)
}

type StateGenerator interface {
	GenerateState(ctx context.Context, ID string) (string, error)
}

type StateVerifier interface {
	VerifyState(ctx context.Context, ID, state string) error
}
