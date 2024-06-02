package adapters

// TODO: Use build tags (i.e. //go:build aws) to conditionally compile this file.

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/oauth2"
)

type DynamoStore struct {
	Client    *dynamodb.Client
	TableName string
	Crypter   TokenCrypter
}

func NewDynamoStore(c context.Context, region string, encryptionKey []byte) (*DynamoStore, error) {
	cfg, cfgErr := config.LoadDefaultConfig(c, func(o *config.LoadOptions) error {
		o.Region = region
		return nil
	})

	cipher, blockErr := NewAesTokenCrypter(encryptionKey)
	if cfgErr != nil || blockErr != nil {
		return nil, errors.Join(cfgErr, blockErr)
	}

	return &DynamoStore{
		Client:    dynamodb.NewFromConfig(cfg),
		TableName: "auth_tokens",
		Crypter:   cipher,
	}, nil
}

func (d *DynamoStore) SaveToken(ctx context.Context, ID string, token *oauth2.Token) error {
	encryptedToken, err := d.Crypter.Encrypt(token)
	if err != nil {
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.TableName,
		Item: map[string]types.AttributeValue{
			"ID":    &types.AttributeValueMemberS{Value: ID},
			"Token": &types.AttributeValueMemberS{Value: encryptedToken},
		},
	})
	return err
}

func (d *DynamoStore) RetrieveToken(ctx context.Context, ID string) (*oauth2.Token, error) {
	result, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.TableName,
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: ID},
		},
	})
	if err != nil {
		return nil, err
	}
	encryptedToken, ok := result.Item["Token"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, errors.New("token not found in record")
	}
	return d.Crypter.Decrypt(encryptedToken.Value)
}
