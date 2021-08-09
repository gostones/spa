package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gostones/spa/internal"
	"github.com/gostones/spa/internal/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cfg internal.Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: internal.ClientVersion,
	Use:     internal.ClientName,
	Short:   "Secure Password Assistant (SPA)",
	Long: `
Secure Password Assistant (SPA)

Unlike password managers, SPA does not store your passwords; it recreates the
same password by taking advantages of the deterministic properties of secure 
cryptographic hashing algorithms. For the same given secret inputs, it always
generates the same output - your password for login with your website and yet
it is infeasible to find out your secret.

SPA is inspired by an ancient technique: split tally. it uses several pieces of
inputs, namely secret, salt, and pepper combined for password generation. It
splits your secret into two overlapping portions: one used for encrypting your
pepper and the other generating passwords.

Secret - that is, a secret, kept to yourself only.

Salt - a piece of text taken from your email, a website, or any other source
that you need to remember to go back if needed. Its hash value, the output of
running secure hash algorithms is stored on your local disk. Salt helps
strengthen your secret - possibly short and weak.

Pepper - an optional piece of data for each webstie to create variations on
the generated passwords. Pepper is encrypted and stored on your local disk.
The pepper should be backed up in the event that you might lose access to your
computer.

A large set of passwords are generated each time. A PIN (password index number)
can be used to pick one of them for your website.

The PIN is the last defense agaist the unlikely worst case senario: your local
computer has been hacked and one of your websites have been compromised by the
same hackers; they have cracked your secret.

SPA can also be used to genearate fake answers to security quesitons required
for password resetting by some websites.
`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true

	rootCmd.PersistentFlags().StringVar(&cfg.BaseDir, "config", "", fmt.Sprintf("custom location for storing salt hash and encrypted pepper (default %s environment variable or $HOME/%s)", spaConfigEnv, defaultDir))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	parseInt := func(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
			return -1
		}
		return v
	}

	if cfg.Secret.Raw == "" {
		cfg.Secret.Raw = os.Getenv(spaSecretEnv)
	}

	if cfg.Pin == -1 {
		cfg.Pin = parseInt(os.Getenv(spaPinEnv))
	}

	if cfg.BaseDir == "" {
		cfg.BaseDir = os.Getenv(spaConfigEnv)
		// or default
		if cfg.BaseDir == "" {
			home, err := homedir.Dir()
			if err != nil {
				exit(err)
			}
			cfg.BaseDir = filepath.Join(home, defaultDir)
		}
	}
	if err := os.MkdirAll(cfg.BaseDir, 0700); err != nil {
		exit(err)
	}
	if err := os.Chmod(cfg.BaseDir, 0700); err != nil {
		exit(err)
	}
}

// exit checks error and exit with the following code:
// 0 -- success - no error
// 1 -- general failure - any standard golang error
// 2 -- usage error
func exit(err error) {
	if err == nil {
		os.Exit(0)
	}

	log.Errorln(err)

	switch err.(type) {
	case *internal.UsageError:
		os.Exit(2)
	default:
		os.Exit(1)
	}
}
