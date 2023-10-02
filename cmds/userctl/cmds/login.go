package cmds

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func GetLoginCommand(root *cli.Root) *cobra.Command {
	var (
		password string
		username string
		totpCode string
		ttl      time.Duration
	)

	cmd := &cobra.Command{
		Use: "login",
		Run: func(cmd *cobra.Command, args []string) {
			if password == "" {
				fmt.Print("Please enter password: ")
				pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()

				if err != nil {
					logrus.Fatalf("failed to read password: %s", err)
				}

				password = string(pwd)
			}

			if username == "" && len(args) > 0 {
				username = args[0]
			}

			if err := root.Login(username, password, totpCode); err != nil {
				logrus.Fatalf("failed to login: %s", err)
			}
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVarP(&password, "password", "p", "", "The password to login")
		flags.StringVarP(&username, "username", "u", "", "The username to login")
		flags.StringVar(&totpCode, "totp-code", "", "The TOTP 2FA code")
		flags.DurationVar(&ttl, "ttl", 0, "The TTL for the access token")
	}

	return cmd
}
