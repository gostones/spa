package cmd

import (
	"github.com/spf13/cobra"
)

// saltCmd represents the salt command
var saltCmd = &cobra.Command{
	DisableFlagsInUseLine: true,
	Use:                   "salt",
	Short:                 "Setup salt",
	Long: `
Salt is a piece of data for generating passwords and answers for use
with a website.

Secret, salt, and pepper are combined to produce passwords and answers to
security questions.

Before you can use SPA, you must setup first by running:

spa salt save --text "<TEXT>"
`,
}

func init() {
	rootCmd.AddCommand(saltCmd)
}
