package cryptoToken

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

type TokenCrypter struct {
	Key string
}

func NewTokenCrypter(key string) *TokenCrypter {
	return &TokenCrypter{Key: key}
}

type EncryptedToken struct {
	IV         string `json:"iv"`
	Ciphertext string `json:"ciphertext"`
	Tag        string `json:"tag"`
}

func GetEncryptedTokenFromBase64(tokenInBase64 string) (EncryptedToken, error) {
	decodedJson, err := base64.StdEncoding.DecodeString(tokenInBase64)
	if err != nil {
		return EncryptedToken{}, err
	}
	var encrypredToken EncryptedToken
	err = json.Unmarshal(decodedJson, &encrypredToken)
	if err != nil {
		return EncryptedToken{}, err
	}
	return encrypredToken, nil
}

func DecryptToken(enc *EncryptedToken, keyBase64 string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	iv, err := base64.StdEncoding.DecodeString(enc.IV)
	if err != nil {
		return "", err
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(enc.Ciphertext)
	if err != nil {
		return "", err
	}

	tagBytes, err := base64.StdEncoding.DecodeString(enc.Tag)
	if err != nil {
		return "", err
	}

	fullCiphertext := append(cipherBytes, tagBytes...)

	plaintext, err := gcm.Open(nil, iv, fullCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
