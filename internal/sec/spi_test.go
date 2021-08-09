package sec

import (
	"bytes"

	"testing"
)

func TestSPIKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}

	secret, _ := RandomBytes(7)
	newSecret, _ := RandomBytes(12)

	// initial key
	key, _ := RandomBytes(3 * 32)
	hash := SPIHash(secret, key)

	newKey := SPIKey(secret, key, newSecret)
	newHash := SPIHash(newSecret, newKey)

	if !bytes.Equal(hash, newHash) {
		t.Fatalf("got: %s want: %s", Base64(hash), Base64(newHash))
	}
	t.Logf("old secret: %s key: %s", Base64(secret), Base64(key))
	t.Logf("new secret: %s key: %s", Base64(newSecret), Base64(newKey))
	t.Logf("hash: %s new: %s", Base64(hash), Base64(newHash))
}

func TestSPIHash(t *testing.T) {
	secret, _ := Debase64("YVvt5kpLyA==")
	key, _ := Debase64("HP+Xu3tHHs3NQMahdAEn0EUkN2kRiZeCBbA7UhzHL5TDZmZZprNEPQEmlP5NYRlBvWDJMAm0E/GoRaFSNmEOBDrXzPKuK4s2jqsCPgLRaxdGAPRkHfgRyutswFeYc3xQ")
	expected, _ := Debase64("IepGF87rOY2+Rkqx2PWVe5tEiNKSOciM1GP+AVfZIHh6J8wCDz97EeX6mI4zZyxyX+kzuI6gZ5kdlA9QHoplrQ==")

	got := SPIHash(secret, key)

	if !bytes.Equal(got, expected) {
		t.Fatalf("got: %s want: %s", Base64(got), Base64(expected))
	}
}
