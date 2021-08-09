package sec

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type cryptData struct {
	key       []byte
	plain     []byte
	cipher    []byte
	iteration int
}

func testCryptData() *cryptData {
	cipher, _ := hex.DecodeString("765d66faad33db3f48f3e674cbeae949d94bba067b842e4cb043992f70499aefa6c26e6524576be4da863cd7aa3ab95f4289d7660ce232951fa3080b0143256df7af12ebcf2866fe")
	return &cryptData{
		plain:     []byte("confidential"),
		key:       []byte("abcde12345"),
		cipher:    cipher,
		iteration: 1,
	}
}

func TestEncrypt(t *testing.T) {
	key := []byte("abcde12345")
	data := []byte("confidential")
	iteration := 1

	t.Logf("key: %v", key)
	ciphertext, err := Encrypt(key, data, iteration)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ciphertext: %s", hex.EncodeToString(ciphertext))

	plaintext, err := Decrypt(key, ciphertext, iteration)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plaintext, data) {
		t.Fatalf("got: %v want: %v", plaintext, data)
	}
	t.Logf("plaintext: %s", plaintext)
}

func TestDecrypt(t *testing.T) {
	tc := testCryptData()
	expected := tc.plain

	dec, err := Decrypt(tc.key, tc.cipher, tc.iteration)
	if !bytes.Equal(dec, expected) {
		t.Fatalf("got: %s want: %s", dec, expected)
	}
	t.Logf("%v %s", err, dec)
}

func BenchmarkEncrypt(b *testing.B) {
	tc := testCryptData()

	for i := 0; i < b.N; i++ {
		Encrypt(tc.key, tc.plain, tc.iteration)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	tc := testCryptData()

	for i := 0; i < b.N; i++ {
		Decrypt(tc.key, tc.cipher, tc.iteration)
	}
}
