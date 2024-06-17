package adapters

import (
	"context"

	"github.com/alexmeuer/juke/pkg/user"
	"github.com/dgraph-io/badger/v4"
	jsoniter "github.com/json-iterator/go"
)

type BadgerStore struct {
	DB *badger.DB
}

func NewBadgerStore(path string) (*BadgerStore, error) {
	if path == "" {
		path = "/tmp/juke_badger_users"
	}
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerStore{
		DB: db,
	}, nil
}

func (b *BadgerStore) Close() error {
	return b.DB.Close()
}

func (b *BadgerStore) CreateUser(ctext context.Context, info user.PublicInfo) error {
	// NOTE: this does not check if the user already exists. It will overwrite.
	return b.DB.Update(func(txn *badger.Txn) error {
		encoded, err := jsoniter.Marshal(info)
		if err != nil {
			return err
		}
		return txn.Set([]byte(info.ID), encoded)
	})
}
