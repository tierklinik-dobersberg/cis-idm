package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"golang.org/x/crypto/ssh/terminal"
)

func getLoginCommand() *cobra.Command {
	passwordAuth := idmv1.PasswordAuth{}

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
			cli := idmv1connect.NewAuthServiceClient(http.DefaultClient, baseURL)

			res, err := cli.Login(context.Background(), connect.NewRequest(&idmv1.LoginRequest{
				AuthType: idmv1.AuthType_AUTH_TYPE_PASSWORD,
				Auth: &idmv1.LoginRequest_Password{
					Password: &passwordAuth,
				},
			}))
			if err != nil {
				return err
			}

			if acm := res.Msg.GetAccessToken(); acm != nil {
				if err := writeTokenFile(acm.Token); err != nil {
					return err
				}

				logrus.Infof("Hello %s", acm.User.DisplayName)
				logrus.Infof("access token stored at %s", tokenPath)

				return nil
			}

			return fmt.Errorf("no access token returned by server")
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVarP(&passwordAuth.Password, "password", "p", "", "The password to login")
		flags.StringVarP(&passwordAuth.Username, "username", "u", "", "The username to login")
	}

	return cmd
}

func writeTokenFile(token string) error {
	return os.WriteFile(tokenPath, []byte(token), 0600)
}
