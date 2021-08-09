package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

const (
	encryptKeyLen = 32
)

func Encrypt(secret, data []byte, iteration int) ([]byte, error) {
	salt := make([]byte, encryptKeyLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key, err := enckey(secret, salt, iteration)
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
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return append(ciphertext, salt...), nil
}

func Decrypt(secret, data []byte, iteration int) ([]byte, error) {
	salt, ciphertext := data[len(data)-encryptKeyLen:], data[:len(data)-encryptKeyLen]
	key, err := enckey(secret, salt, iteration)
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
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func enckey(secret, salt []byte, iteration int) ([]byte, error) {
	key, err := SPA(secret, salt, encryptKeyLen, iteration)
	if err != nil {
		return nil, err
	}

	return key, nil
}
