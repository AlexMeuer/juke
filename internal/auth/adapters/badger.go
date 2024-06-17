package adapters

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const stateLength = 64

type BadgerStore struct {
	DB      *badger.DB
	Crypter TokenCrypter
}

func NewBadgerStore(path string, encryptionKey []byte) (*BadgerStore, error) {
	if path == "" {
		path = "/tmp/juke_badger_tkns"
	}
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	cipher, err := NewAesTokenCrypter(encryptionKey)
	if err != nil {
		return nil, err
	}

	return &BadgerStore{
		DB:      db,
		Crypter: cipher,
	}, nil
}

func (b *BadgerStore) Close() error {
	return b.DB.Close()
}

func (b *BadgerStore) SaveToken(c context.Context, ID string, token *oauth2.Token) error {
	encryptedToken, err := b.Crypter.Encrypt(token, []byte(ID))
	if err != nil {
		return err
	}

	return b.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(ID), []byte(encryptedToken))
	})
}

func (b *BadgerStore) RetrieveToken(c context.Context, ID string) (*oauth2.Token, error) {
	var token *oauth2.Token
	IDBytes := []byte(ID)

	err := b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(IDBytes)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			token, err = b.Crypter.Decrypt(string(val), IDBytes)
			return err
		})
	})

	return token, err
}

func (b *BadgerStore) GenerateState(c context.Context, ID string) (string, error) {
	buffer := make([]byte, stateLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	state := base64.URLEncoding.EncodeToString(buffer)

	err = b.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(ID), []byte(state))
	})
	return state, err
}

func (b *BadgerStore) VerifyState(c context.Context, ID, state string) error {
	return b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			if string(val) != state {
				log.Error().
					Str("stored_state", string(val)).
					Str("state", state).
					Msg("state mismatch")
				return errors.New("state mismatch")
			}
			return nil
		})
	})
}
