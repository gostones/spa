package cmd

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/gostones/spa/internal/log"
	"github.com/gostones/spa/internal/sec"
)

// normalize strips all spaces
func normalize(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func checkSaltSecret() error {
	if cfg.Salt.Raw != "" {
		salt, err := hashSalt(normalizedSalt())
		if err != nil {
			return err
		}
		cfg.Salt.Hash = salt
	} else {
		p := saltFilename()
		if !checkFile(p) {
			return fmt.Errorf("salt is required. please run 'spa salt -h' for details")
		}
		hash, err := readSaltFile(p)
		if err != nil {
			return err
		}
		cfg.Salt.Hash = hash
	}

	// TODO this has to run after salt hash
	if cfg.Secret.Raw == "" {
		raw, err := log.PromptSecret(secretPrompt)
		if err != nil {
			return err
		}
		cfg.Secret.Raw = raw
	}

	secrets, err := hashSecret(cfg.Secret.Raw, cfg.Salt.Hash)
	if err != nil {
		return err
	}
	cfg.Secret.Stock = secrets[0]
	cfg.Secret.Foil = secrets[1]

	return nil
}

func generator(codebook string) (func(string, string, string, int) ([]string, error), error) {
	g := sec.KeyGen(codebook, cfg.Secret.Stock, cfg.Salt.Hash, hashKeyLen, keyGenIteration)
	if g == nil {
		return nil, fmt.Errorf("failed to create password generator")
	}
	return g, nil
}

func spakdf(key, salt []byte, iteration int) ([]byte, error) {
	ba, err := sec.SPAKDF(key, salt, masterKeyCount, hashKeyLen, iteration)
	if err != nil {
		return nil, err
	}
	var h []byte
	for _, k := range ba {
		h = append(h, k...)
	}
	return h[0:], nil
}

func hashSalt(raw []byte) ([]byte, error) {
	key, salt := raw[2:], raw[0:2]
	return spakdf(key, salt, saltIteration)
}

func hashSecret(raw string, salt []byte) ([][]byte, error) {
	if len(raw) < minSecretLen {
		return nil, ErrSecretTooShort
	}

	key, err := readKey()
	if err != nil {
		return nil, err
	}
	hash := sec.SPIHash([]byte(raw), key)
	// split so we have two secrets:
	// one for encryption and the other for password generation
	secrets := split2(hash)

	var wg sync.WaitGroup

	var keys [2][]byte
	var errs [2]error

	spa := func(i int, key []byte) {
		defer wg.Done()
		idx := sec.FNV(key, uint32(len(salt)-hashKeyLen))
		keys[i], errs[i] = spakdf(key, salt[idx:idx+hashKeyLen], secretIteration)
	}

	for i, v := range secrets {
		wg.Add(1)
		go spa(i, v)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return keys[0:], nil
}

func readSaltFile(p string) ([]byte, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	salt, err := hex.DecodeString(string(b))
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func catSaltFile(p string) error {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	io.WriteString(os.Stdout, string(b))
	io.WriteString(os.Stdout, "\n")
	return nil
}

func saveSaltFile(p string, salt []byte) error {
	perm := os.FileMode(0600)
	s := hex.EncodeToString(salt)
	if err := ioutil.WriteFile(p, []byte(s), perm); err != nil {
		return err
	}
	return os.Chmod(p, perm)
}

func checkFile(p string) bool {
	var _, err = os.Stat(p)
	return !os.IsNotExist(err)
}

func split2(input []byte) [][]byte {
	n := len(input)
	x := (n * 3) / 10
	return [][]byte{input[:n-x], input[x:]}
}

func enterNewPassword() (string, error) {
	raw, err := log.PromptSecret(secretNewPrompt)
	if err != nil {
		return "", err
	}
	if len(raw) < minSecretLen {
		return "", ErrSecretTooShort
	}

	rawAgain, err := log.PromptSecret(secretNewAgainPrompt)
	if err != nil {
		return "", err
	}
	if raw != rawAgain {
		return "", ErrSecretMismatch
	}

	return raw, nil
}

func requireSecret() error {
	if err := checkSaltSecret(); err != nil {
		exit(err)
	}
	// validate secret
	_, err := decryptSafe(cfg.Secret.Foil)
	return err
}
