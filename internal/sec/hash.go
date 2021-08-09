package sec

import (
	"crypto/hmac"
	"crypto/sha512"
	"hash/fnv"
	"io"
	"strconv"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#work-factors

type phaFunc func([]byte, []byte, int) ([]byte, error)

// list of password hashing algorithms
var algos = []phaFunc{
	scryptKey,
	pbkdf2Key,
	argon2idKey,
}

func SPA(secret, salt []byte, keyLen, iteration int) ([]byte, error) {
	if iteration < 1 {
		iteration = 1
	}

	pha := func(v int64) phaFunc {
		cnt := len(algos)
		which := v % int64(cnt)
		return algos[which]
	}

	// pick a permutation based on salt
	sch, err := schedule(int64(iteration), salt)
	if err != nil {
		return nil, err
	}

	// derive salt for each iteration
	jar, err := KDF(salt, iteration, keyLen)
	if err != nil {
		return nil, err
	}

	hash := secret
	for i, v := range sch {
		hash, err = pha(v)(hash, jar[i], keyLen)
		if err != nil {
			return nil, err
		}
	}
	return hash, nil
}

func pbkdf2Key(pwd, salt []byte, keyLen int) ([]byte, error) {
	iteration := 120000
	h := pbkdf2.Key(pwd, salt, iteration, keyLen, sha512.New)
	return h, nil
}

func argon2idKey(pwd, salt []byte, keyLen int) ([]byte, error) {
	time := 1
	memory := 64 * 1024
	threads := 4
	return argon2.IDKey(pwd, salt, uint32(time), uint32(memory), uint8(threads), uint32(keyLen)), nil
}

func scryptKey(pwd, salt []byte, keyLen int) ([]byte, error) {
	N := 32768
	r := 8
	p := 1
	return scrypt.Key(pwd, salt, N, r, p, keyLen)
}

func KDF(secret []byte, keyCount, keyLen int) ([][]byte, error) {
	kdf := func(salt, info []byte, count, klen int) ([][]byte, error) {
		hash := sha512.New
		hf := hkdf.New(hash, secret, salt, info)
		var keys [][]byte
		for i := 0; i < count; i++ {
			key := make([]byte, klen)
			if _, err := io.ReadFull(hf, key); err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		return keys, nil
	}

	var h [][]byte
	for i := 0; i < keyCount; i++ {
		k, err := kdf([]byte(strconv.Itoa(i)), nil, 1, keyLen)
		if err != nil {
			return nil, err
		}
		h = append(h, k[0])
	}
	return h, nil
}

func SPAKDF(raw, salt []byte, keyCount, keyLen, iteration int) ([][]byte, error) {
	data, err := SPA(raw, salt, keyLen, iteration)
	if err != nil {
		return nil, err
	}
	ba, err := KDF(data, keyCount, keyLen)
	if err != nil {
		return nil, err
	}
	return ba, nil
}

func HMAC(key, data []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func FNV(b []byte, n uint32) uint32 {
	f := fnv.New32a()
	f.Write(b)
	return f.Sum32() % uint32(n)
}

func HMACFNV(x, y []byte, max uint32) uint32 {
	b := HMAC(x, y)
	return FNV(b, max)
}
