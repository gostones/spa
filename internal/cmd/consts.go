package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// total number of websites: 1.88 billion
	// < 2^16*2^16
	hashKeyLen     = 64
	masterKeyCount = 1024

	saltIteration   = 1024
	secretIteration = 8
	keyGenIteration = 8
	cryptIteration  = 6

	// 1 year < 2^26 seconds
	// total number of pc < 2^34 ( every one owns 2 devices 2^33 * 2)
	// letter (2^6) vs byte (2^8)
	minSaltLen   = 86
	minSecretLen = 6

	autoPepperSize = 64

	minPwdLength     = 8
	maxPwdLength     = 64
	defaultPwdLength = 32

	defaultMaxPIN = 10
)

const (
	spaSecretEnv = "SPA_SECRET"
	spaPinEnv    = "SPA_PIN"

	spaConfigEnv = "SPA_CONFIG"

	spaSaltFileEnv   = "SPA_SALT_FILE"
	spaPepperFileEnv = "SPA_PEPPER_FILE"
	spaKeyFileEnv    = "SPA_KEY_FILE"
)

const defaultDir = ".spa"

const (
	secretPrompt = `Enter your secret here: `

	secretNewPrompt      = `Enter your new secret here: `
	secretNewAgainPrompt = `Enter your new secret again: `

	autoPepperGeneration = `Generating random pepper...
`
	savePepperOverrite = `
Warning: you are about to change pepper for %q

from: %q
to: %q

The new value will be used for password generation.

You should consider backing up your pepper in case you lose access to this computer. 

`
	savePepperPrompt = `Continue? [y/N] `
)

var ErrSecretTooShort = fmt.Errorf("secret is too short. minimum characters required: %v", minSecretLen)
var ErrSecretMismatch = fmt.Errorf("secret does not match! Please try again")

func pepperFilename() string {
	name := filepath.Join(cfg.BaseDir, "pepper")
	if envName := os.Getenv(spaPepperFileEnv); envName != "" {
		name = envName
	}
	return name
}

func keyFilename() string {
	name := filepath.Join(cfg.BaseDir, "key")
	if envName := os.Getenv(spaKeyFileEnv); envName != "" {
		name = envName
	}
	return name
}

func saltFilename() string {
	name := filepath.Join(cfg.BaseDir, "salt")
	if envName := os.Getenv(spaSaltFileEnv); envName != "" {
		name = envName
	}
	return name
}
