package sec

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func S2b(s string) []byte {
	return []byte(s)
}

func B2s(b []byte) string {
	return string(b)
}

func Base64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Debase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
