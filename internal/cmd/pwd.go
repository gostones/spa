package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gostones/spa/internal"
	"github.com/gostones/spa/internal/log"
	"github.com/gostones/spa/internal/sec"
)

func genPwd(cmd *cobra.Command) error {
	peppers, err := readSite(cmd)
	if err != nil {
		return err
	}

	c := cfg.Pwd
	count := cfg.Count
	if cfg.Pin >= 0 {
		count = cfg.Pin + 1
	}

	codebook := sec.MakeCodebook(sec.AlphaNumericSymbol, cfg.Pwd.Mask)
	g, err := generator(codebook)
	if err != nil {
		return err
	}

	pwds, err := g(cfg.Domain, cfg.User, c.Pepper, count)
	if err != nil {
		return err
	}

	n := c.Length
	print := func(pin int) {
		s := fmt.Sprintf("[%04v] %s", pin, pwds[pin][0:n])
		log.Infoln(s)
	}

	if cfg.Pin < 0 {
		for i := 0; i < len(pwds); i++ {
			print(i)
		}
	} else {
		print(cfg.Pin)
	}

	if err := saveSite(peppers); err != nil {
		return err
	}

	return nil
}

func readSite(cmd *cobra.Command) (map[string]internal.PwdConfig, error) {
	peppers, err := readPepper(cfg.Domain, cfg.User)
	if err != nil {
		return nil, err
	}
	du := domainUser(cfg.Domain, cfg.User)
	c := peppers[du]
	// default/old values if no new values are provided on command line.
	if !cmd.Flags().Changed("pepper") {
		cfg.Pwd.Pepper = c.Pepper
	}
	if !cmd.Flags().Changed("mask") {
		cfg.Pwd.Mask = c.Mask
	}
	if !cmd.Flags().Changed("length") {
		cfg.Pwd.Length = c.Length
	}

	//
	if cmd.Flags().Changed("pepper") {
		pepper := cfg.Pwd.Pepper
		if pepper == "auto" {
			log.Infof(autoPepperGeneration)

			b, err := sec.RandomBytes(autoPepperSize)
			if err != nil {
				return nil, err
			}
			pepper = sec.Base64(b)
		}

		log.Infof(savePepperOverrite, cfg.Domain, c.Pepper, pepper)
		choice, err := log.Confirm(savePepperPrompt)
		if err != nil {
			return nil, err
		}
		switch choice {
		case "y":
			cfg.Pwd.Pepper = pepper
		case "n":
			return nil, fmt.Errorf("")
		}
	}

	// update
	peppers[du] = cfg.Pwd

	return peppers, nil
}

func saveSite(peppers map[string]internal.PwdConfig) error {
	return writePepper(peppers)
}

func validatePwdFlags(cmd *cobra.Command, args []string) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain name is required")
	}

	if cfg.Pwd.Length < minPwdLength || cfg.Pwd.Length > maxPwdLength {
		return fmt.Errorf("invalid length: %v. valid range [%v, %v]", cfg.Pwd.Length, minPwdLength, maxPwdLength)
	}

	return nil
}

// pwdCmd represents the pwd command
var pwdCmd = &cobra.Command{
	DisableFlagsInUseLine: true,
	Use:                   "pwd -d <DOMAIN NAME> [-p <PIN>] [--pepper auto]",
	Short:                 "Generate passwords",
	Long: `
Generate a set of candidate passwords to be used.

A number of passwords will be printed, each password is prepended with a
password index number (PIN).

You can pick any password for use with your web site. Just remember to use
the same PIN for the same site.
`,
	Args: validatePwdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkSaltSecret(); err != nil {
			exit(err)
		}

		err := genPwd(cmd)
		exit(err)
	},
}

func init() {
	rootCmd.AddCommand(pwdCmd)

	pwdCmd.Flags().StringVar(&cfg.Salt.Raw, "salt", "", "specify a salt to use, whitespaces are ignored. default: saved in .spa/salt")
	pwdCmd.Flags().StringVarP(&cfg.Secret.Raw, "secret", "s", "", "your secret.")

	pwdCmd.Flags().VarP(newDomainValue("", &cfg.Domain), "domain", "d", "domain name of the web site. case insensitive. e.g. example.com")
	pwdCmd.Flags().StringVarP(&cfg.User, "user", "u", "", "optional username, email, or ID for the web site. provide a value if you have multiple accounts with the same web site.")
	pwdCmd.Flags().StringVar(&cfg.Pwd.Pepper, "pepper", "", "recommended text to generate a different password. enter 'auto' if you want it randomly generated.")
	pwdCmd.Flags().StringVar(&cfg.Pwd.Note, "note", "", "attach a note for information only, note does not alter password generation.")

	pwdCmd.Flags().StringVar(&cfg.Pwd.Mask, "mask", sec.EncloseEscape, "optional characters to exclude from generated password. you should only provide a mask for web sites that do not accept certain special symbols.")
	pwdCmd.Flags().IntVar(&cfg.Pwd.Length, "length", defaultPwdLength, fmt.Sprintf("optional password length, minimum %v maximum %v. it is recommended you use the full length unless the web site sets a limit.", minPwdLength, maxPwdLength))

	pwdCmd.Flags().IntVar(&cfg.Count, "count", defaultMaxPIN, "optional number of passwords to generate")
	pwdCmd.Flags().VarP(newPinValue(-1, &cfg.Pin), "pin", "p", "optional number to pick the password, the full list will be shown if not provided.")

	pwdCmd.MarkFlagRequired("domain")

	pwdCmd.Flags().MarkHidden("salt")
	pwdCmd.Flags().MarkHidden("secret")
}
