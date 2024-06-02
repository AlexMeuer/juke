package adapters

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type TokenCrypter interface {
	Encrypt(token *oauth2.Token, IV []byte) (string, error)
	Decrypt(encryptedToken string, IV []byte) (*oauth2.Token, error)
}

type BlockCipherTokenCrypter struct {
	cipher.Block
}

func NewAesTokenCrypter(key []byte) (*BlockCipherTokenCrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &BlockCipherTokenCrypter{block}, nil
}

func (c *BlockCipherTokenCrypter) Encrypt(token *oauth2.Token, IV []byte) (string, error) {
	plaintext, err := jsoniter.Marshal(token)
	if err != nil {
		return "", err
	}
	// FIXME: (Unit test this) What if IV is < BlockSize? What about other cases?
	cfb := cipher.NewCFBEncrypter(c.Block, IV[:c.BlockSize()])
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *BlockCipherTokenCrypter) Decrypt(encryptedToken string, IV []byte) (*oauth2.Token, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBDecrypter(c.Block, IV[:c.BlockSize()])
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	var token oauth2.Token
	if err := jsoniter.Unmarshal(plaintext, &token); err != nil {
		return nil, err
	}
	return &token, nil
}
