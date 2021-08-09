package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gostones/spa/internal"
	"github.com/gostones/spa/internal/log"
	"github.com/gostones/spa/internal/sec"
)

type Safe struct {
	Nonce []byte
	Data  map[string]internal.PwdConfig
}

func domainUser(domain, user string) string {
	return strings.Join([]string{domain, user}, ":")
}

func readPepper(domain, user string) (map[string]internal.PwdConfig, error) {
	s, err := decryptSafe(cfg.Secret.Foil)
	if err != nil {
		return nil, err
	}
	peppers := s.Data

	du := domainUser(domain, user)
	if _, ok := peppers[du]; !ok {
		b, err := sec.RandomBytes(autoPepperSize)
		if err != nil {
			return nil, err
		}
		c := internal.PwdConfig{
			Pepper: sec.Base64(b),
			Length: defaultPwdLength,
			Mask:   sec.EncloseEscape,
		}
		peppers[du] = c
	}

	return peppers, nil
}

func writePepper(peppers map[string]internal.PwdConfig) error {
	s, err := decryptSafe(cfg.Secret.Foil)
	if err != nil {
		return err
	}
	s.Data = peppers

	return encryptSafe(cfg.Secret.Foil, s)
}

func pickKey(key, nonce []byte) []byte {
	hi := func(x, y []byte, max int) int {
		b := sec.HMAC(x, y)
		return int(sec.FNV(b, uint32(max)))
	}
	idx := hi(key[0:hashKeyLen], nonce, len(key)-hashKeyLen)
	secret := key[idx : idx+hashKeyLen]
	return secret
}

func encryptSafe(key []byte, s *Safe) error {
	nonce, err := sec.RandomBytes(hashKeyLen)
	if err != nil {
		return err
	}
	s.Nonce = nonce

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	secret := pickKey(key, nonce)
	cipher, err := sec.Encrypt(secret, data, cryptIteration)
	if err != nil {
		return err
	}
	enc := base64.StdEncoding.EncodeToString(append(cipher, nonce...))

	file := pepperFilename()
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	perm := os.FileMode(0600)
	if err := ioutil.WriteFile(file, []byte(enc), perm); err != nil {
		return err
	}
	return os.Chmod(file, perm)
}

func decryptSafe(key []byte) (*Safe, error) {
	file := pepperFilename()
	if !checkFile(file) {
		return &Safe{
			Data: make(map[string]internal.PwdConfig),
		}, nil
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dec, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, err
	}

	nonce := dec[len(dec)-hashKeyLen:]
	cipher := dec[0 : len(dec)-hashKeyLen]
	secret := pickKey(key, nonce)
	data, err := sec.Decrypt(secret, cipher, cryptIteration)
	if err != nil {
		return nil, err
	}

	var s Safe
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func catPepperFile() error {
	s, err := decryptSafe(cfg.Secret.Foil)
	if err != nil {
		return err
	}
	peppers := s.Data

	pretty := func(c internal.PwdConfig) ([]byte, error) {
		b, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		var out bytes.Buffer
		if err := json.Indent(&out, b, "", "  "); err != nil {
			return nil, err
		}
		return out.Bytes(), err
	}

	var keys []string
	for k := range peppers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	domain := strings.ToLower(cfg.Domain)
	user := strings.ToLower(cfg.User)
	filtering := (domain != "" || user != "")
	match := func(s string) bool {
		sa := strings.Split(s, ":")
		return (domain == "" || strings.Contains(strings.ToLower(sa[0]), domain)) && (user == "" || strings.Contains(strings.ToLower(sa[1]), user))
	}

	for _, k := range keys {
		if filtering && !match(k) {
			continue
		}
		v := peppers[k]
		b, err := pretty(v)
		if err != nil {
			return nil
		}
		log.Infof("%q: %s\n", k, string(b))
	}
	return nil
}

// readKey reads the key; randmoly generates one if not found.
func readKey() ([]byte, error) {
	file := keyFilename()

	if !checkFile(file) {
		b, err := sec.InitSPIKey()
		if err != nil {
			return nil, err
		}
		return b, writeKey(b)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(b))
}

func writeKey(newKey []byte) error {
	file := keyFilename()

	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	enc := base64.StdEncoding.EncodeToString(newKey)
	perm := os.FileMode(0600)
	if err := ioutil.WriteFile(file, []byte(enc), perm); err != nil {
		return err
	}
	return os.Chmod(file, perm)
}
