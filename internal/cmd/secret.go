package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gostones/spa/internal/sec"
)

func changeSecret() error {
	if cfg.Secret.NewRaw == "" {
		raw, err := enterNewPassword()
		if err != nil {
			return err
		}
		cfg.Secret.NewRaw = raw
	}

	if len(cfg.Secret.NewRaw) < minSecretLen {
		return ErrSecretTooShort
	}

	key, err := readKey()
	if err != nil {
		return err
	}
	newKey := sec.SPIKey([]byte(cfg.Secret.Raw), key, []byte(cfg.Secret.NewRaw))
	return writeKey(newKey)
}

// secretCmd represents the server command
var secretCmd = &cobra.Command{
	Hidden:                false,
	DisableFlagsInUseLine: true,
	Use:                   "secret",
	Short:                 "Change secret",
	Long: `
Change your secret or set a new one (if first time).

A key file will be generated or updated. default: ~/.spa/key. You should back up this file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// verify old secret if key file exists
		if checkFile(keyFilename()) {
			if err := requireSecret(); err != nil {
				exit(err)
			}
		}
		err := changeSecret()
		exit(err)
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.Flags().StringVarP(&cfg.Secret.Raw, "secret", "s", "", "your old secret")
	secretCmd.Flags().StringVarP(&cfg.Secret.NewRaw, "new-secret", "n", "", "your new secret")

	secretCmd.Flags().MarkHidden("secret")
	secretCmd.Flags().MarkHidden("new-secret")
}
