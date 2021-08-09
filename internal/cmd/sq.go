package cmd

import (
	"fmt"

	"github.com/gostones/spa/internal/log"
	"github.com/gostones/spa/internal/sec"
	"github.com/spf13/cobra"
)

func validateSecureQuestionFlags(cmd *cobra.Command, args []string) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain name is required")
	}

	if cfg.Question.Question == "" {
		return fmt.Errorf("security question is required")
	}

	return nil
}

func genSQ() error {
	codebook := sec.MakeCodebook(sec.AlphaNumeric, "")
	g, err := generator(codebook)
	if err != nil {
		return err
	}

	c := cfg.Question

	count := cfg.Count
	if cfg.Pin >= 0 {
		count = cfg.Pin + 1
	}
	//
	answers, err := g(cfg.Domain, cfg.User, normalize(c.Question), count)
	if err != nil {
		return err
	}

	min := 16
	print := func(pin int) {
		answer := answers[pin]
		// only part of the answer
		n := min + int(sec.FNV([]byte(answer), uint32(min)))
		s := fmt.Sprintf("[%04v] %s", pin, sec.SpaceOut(answer[0:n], 0))
		log.Infoln(s)
	}

	if cfg.Pin < 0 {
		for i := 0; i < len(answers); i++ {
			print(i)
		}
	} else {
		print(cfg.Pin)
	}

	return nil
}

// sqCmd represents the sq command
var sqCmd = &cobra.Command{
	DisableFlagsInUseLine: true,
	Use:                   "sq -d <DOMAIN NAME> -q <SECURITY QUESTION> [-p <PIN>]",
	Short:                 "Generate fake answers",
	Long: `
Generate fake answers to security questions.
`,
	Args: validateSecureQuestionFlags,
	Run: func(cmd *cobra.Command, args []string) {
		if err := requireSecret(); err != nil {
			exit(err)
		}

		err := genSQ()
		exit(err)
	},
}

func init() {
	rootCmd.AddCommand(sqCmd)

	sqCmd.Flags().StringVar(&cfg.Salt.Raw, "salt", "", "specify a salt to use, whitespaces are ignored. default: saved in ~/.spa/salt")
	sqCmd.Flags().StringVarP(&cfg.Secret.Raw, "secret", "s", "", "your secret")

	sqCmd.Flags().VarP(newDomainValue("", &cfg.Domain), "domain", "d", "domain name of the web site. case insensitive. e.g. example.com.")
	sqCmd.Flags().StringVarP(&cfg.User, "user", "u", "", "optional username, email, or ID for the web site. provide a value if you have multiple accounts with the web site.")
	sqCmd.Flags().StringVarP(&cfg.Question.Question, "question", "q", "", "security question, whitespaces are ignored.")

	sqCmd.Flags().IntVar(&cfg.Count, "count", defaultMaxPIN, "optional number of answers to generate")
	sqCmd.Flags().VarP(newPinValue(-1, &cfg.Pin), "pin", "p", "optional number to pick the answer")

	sqCmd.MarkFlagRequired("domain")
	sqCmd.MarkFlagRequired("question")

	sqCmd.Flags().MarkHidden("secret")
	sqCmd.Flags().MarkHidden("salt")
}
