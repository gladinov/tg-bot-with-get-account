package cryptoToken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type EncryptedToken struct {
	IV         string `json:"iv"`
	Ciphertext string `json:"ciphertext"`
	Tag        string `json:"tag"`
}

func (e *EncryptedToken) ToBase64() (string, error) {
	const op = "cryptoToken.ToBase64"
	jsonData, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}
	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return encoded, nil
}

type TokenCrypter struct {
	Key string
}

func NewTokenCrypter(key string) *TokenCrypter {
	return &TokenCrypter{Key: key}
}

func EncryptToken(plaintext string, keyBase64 string) (*EncryptedToken, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, iv, []byte(plaintext), nil)

	return &EncryptedToken{
		IV:         base64.StdEncoding.EncodeToString(iv),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext[:len(ciphertext)-gcm.Overhead()]),
		Tag:        base64.StdEncoding.EncodeToString(ciphertext[len(ciphertext)-gcm.Overhead():]),
	}, nil
}

func NewAES(keyB64 string) (cipher.Block, error) {
	// decode base64 key
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}

	// must be exactly 32 bytes
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: %d bytes (expected 32)", len(key))
	}

	return aes.NewCipher(key)
}
