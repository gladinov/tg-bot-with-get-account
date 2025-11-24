package cryptoToken

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type EncryptedToken struct {
	IV         string `json:"iv"`
	Ciphertext string `json:"ciphertext"`
	Tag        string `json:"tag"`
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
