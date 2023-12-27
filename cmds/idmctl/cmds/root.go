package cmds

import (
	"encoding/json"
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
)

func PrepareRootCommand(root *cli.Root) {
	root.AddCommand(
		GetLoginCommand(root),
		GetProfileCommand(root),
		GetUsersCommand(root),
		GetRegisterUserCommand(root),
		GetRoleCommand(root),
		GetSendNotificationCommand(root),
		GenerateVAPIDKeys(),
	)
}

func GenerateVAPIDKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "gen-vapid-keys",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			prv, pub, err := webpush.GenerateVAPIDKeys()
			if err != nil {
				logrus.Fatalf("error: %s", err)
			}

			cfg := config.Config{
				WebPush: &config.WebPush{
					VAPIDpublicKey:  pub,
					VAPIDprivateKey: prv,
				},
			}

			blob, err := json.Marshal(cfg)
			if err != nil {
				logrus.Fatal(err)
			}

			blob, err = yaml.JSONToYAML(blob)
			if err != nil {
				logrus.Fatal(err)
			}

			fmt.Println(string(blob))
		},
	}

	return cmd
}
