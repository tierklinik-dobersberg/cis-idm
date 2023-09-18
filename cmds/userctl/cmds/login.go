package cmds

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/protobuf/types/known/durationpb"
)

func GetLoginCommand(root *cli.Root) *cobra.Command {
	passwordAuth := idmv1.PasswordAuth{}
	var totpCode string
	var ttl time.Duration

	cmd := &cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			if passwordAuth.Password == "" {
				fmt.Print("Please enter password: ")
				pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()

				if err != nil {
					return err
				}

				passwordAuth.Password = string(pwd)
			}

			// we use http.DefaultClient here so we don't get an error if the access token
			// provided by the custom httpClient transport is invalid
			cli := idmv1connect.NewAuthServiceClient((&http.Client{Transport: root.Transport}), root.BaseURLS.Idm)

			var ttlpb *durationpb.Duration
			if ttl > 0 {
				ttlpb = durationpb.New(ttl)
			}

			res, err := cli.Login(context.Background(), connect.NewRequest(&idmv1.LoginRequest{
				AuthType: idmv1.AuthType_AUTH_TYPE_PASSWORD,
				Auth: &idmv1.LoginRequest_Password{
					Password: &passwordAuth,
				},
				Ttl: ttlpb,
			}))
			if err != nil {
				return err
			}

			// check if we need to handle 2fa
			if mfa := res.Msg.GetMfaRequired(); mfa != nil {
				switch mfa.Kind {
				case idmv1.RequiredMFAKind_REQUIRED_MFA_KIND_TOTP:
					if totpCode == "" {
						fmt.Print("Please enter TOTP code: ")
						code, err := terminal.ReadPassword(int(os.Stdin.Fd()))
						fmt.Println()

						if err != nil {
							logrus.Fatal(err)
						}

						totpCode = string(code)
					}

					res, err = cli.Login(context.Background(), connect.NewRequest(&idmv1.LoginRequest{
						AuthType: idmv1.AuthType_AUTH_TYPE_TOTP,
						Auth: &idmv1.LoginRequest_Totp{
							Totp: &idmv1.TotpAuth{
								Code:  totpCode,
								State: mfa.State,
							},
						},
						Ttl: ttlpb,
					}))

					if err != nil {
						logrus.Fatal(err)
					}

				default:
					logrus.Fatalf("unsupported mfa kind required: %s", mfa.Kind.String())
				}
			}

			if acm := res.Msg.GetAccessToken(); acm != nil {
				if err := writeTokenFile(root.TokenPath, acm.Token); err != nil {
					return err
				}

				logrus.Infof("Hello %s", acm.User.DisplayName)
				logrus.Infof("access token stored at %s", root.TokenPath)

				return nil
			}

			return fmt.Errorf("no access token returned by server")
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVarP(&passwordAuth.Password, "password", "p", "", "The password to login")
		flags.StringVarP(&passwordAuth.Username, "username", "u", "", "The username to login")
		flags.StringVar(&totpCode, "totp-code", "", "The TOTP 2FA code")
		flags.DurationVar(&ttl, "ttl", 0, "The TTL for the access token")
	}

	return cmd
}

func writeTokenFile(tokenPath string, token string) error {
	return os.WriteFile(tokenPath, []byte(token), 0600)
}
