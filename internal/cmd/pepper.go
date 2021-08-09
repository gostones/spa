package cmd

import (
	"github.com/spf13/cobra"
)

// pepperCmd represents the server command
var pepperCmd = &cobra.Command{
	Hidden:                false,
	DisableFlagsInUseLine: true,
	Use:                   "pepper",
	Short:                 "Show pepper",
	Long: `
Print the content of pepper to the standard output.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkSaltSecret(); err != nil {
			exit(err)
		}

		err := catPepperFile()
		exit(err)
	},
}

func init() {
	rootCmd.AddCommand(pepperCmd)

	pepperCmd.Flags().StringVarP(&cfg.Domain, "domain", "d", "", "filter by domain name")
	pepperCmd.Flags().StringVarP(&cfg.User, "user", "u", "", "filter by user")

	pepperCmd.Flags().StringVarP(&cfg.Secret.Raw, "secret", "s", "", "your secret")

	pepperCmd.Flags().MarkHidden("secret")
}
