package cmd

import (
	"github.com/spf13/cobra"
)

func showSalt() error {
	p := saltFilename()
	if checkFile(p) {
		return catSaltFile(p)
	}
	return nil
}

// saltCatCmd represents the cat command
var saltCatCmd = &cobra.Command{
	DisableFlagsInUseLine: true,
	Use:                   "cat",
	Short:                 "Print salt",
	Long: `
Print the saved hash value of your salt to the standard output.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := showSalt()
		exit(err)
	},
}

func init() {
	saltCmd.AddCommand(saltCatCmd)
}
