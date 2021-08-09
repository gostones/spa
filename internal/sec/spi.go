package sec

import (
	"encoding/binary"
)

const spiMax uint32 = 1<<16 - 1

// SPIKey returns a new key for the new secret given the old secret and key,
// the new secret and key will hash to the same value with hmac/fnv.
// The old key is divided into chunks of 3 bytes and used for the computation
// of the new key based on hmac/fnv hash.
func SPIKey(secret, key, newSecret []byte) []byte {
	collision := func(n uint32) []byte {
		for i := 0; ; i++ {
			k, _ := RandomBytes(3)
			if HMACFNV(newSecret, k, spiMax) == n {
				return k
			}
		}
	}
	var found []byte
	for i := 0; i < len(key); i += 3 {
		// invariant hash
		h := HMACFNV(secret, key[i:i+3], spiMax)
		h %= spiMax
		c := collision(h)
		found = append(found, c...)
	}
	return found
}

func SPIHash(secret, key []byte) []byte {
	size := 2 * (len(key) / 3)
	b := make([]byte, size+2)

	for i, j := 0, 0; i < len(key); i, j = i+3, j+2 {
		n := HMACFNV(secret, key[i:i+3], spiMax)
		binary.LittleEndian.PutUint32(b[j:], n)
	}
	return b[0:size]
}

// InitSPIKey return a random initial key
func InitSPIKey() ([]byte, error) {
	return RandomBytes(3 * 32)
}
