package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"golang.org/x/crypto/ssh/terminal"
)

func getProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "profile [sub-command]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := idmv1connect.NewAuthServiceClient(httpClient, baseURL)

			res, err := cli.Introspect(context.Background(), connect.NewRequest(&idmv1.IntrospectRequest{}))
			if err != nil {
				return err
			}

			if profile := res.Msg.Profile; profile != nil {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")

				return enc.Encode(profile)
			}

			return nil
		},
	}

	cmd.AddCommand(
		getChangePasswordCommand(),
		getAddEmailCommand(),
		getDeleteEmailCommand(),
	)

	return cmd
}

func getChangePasswordCommand() *cobra.Command {
	req := idmv1.ChangePasswordRequest{}

	cmd := &cobra.Command{
		Use: "set-password",
		RunE: func(cmd *cobra.Command, args []string) error {
			if req.OldPassword == "" {
				fmt.Print("Please enter current password: ")
				pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					return err
				}

				req.OldPassword = string(pwd)
			}

			if req.NewPassword == "" {
				fmt.Print("Please enter new password: ")
				pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					return err
				}

				req.NewPassword = string(pwd)
			}

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			_, err := cli.ChangePassword(context.Background(), connect.NewRequest(&req))
			if err != nil {
				return err
			}

			return nil
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVarP(&req.NewPassword, "new-password", "n", "", "The new user password")
		flags.StringVarP(&req.OldPassword, "password", "p", "", "The old user password")
	}

	return cmd
}

func getAddEmailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-mail",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := &idmv1.AddEmailAddressRequest{
				Email: args[0],
			}

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)
			_, err := cli.AddEmailAddress(context.Background(), connect.NewRequest(msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func getDeleteEmailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete-mail",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := &idmv1.DeleteEmailAddressRequest{
				Id: args[0],
			}

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)
			_, err := cli.DeleteEmailAddress(context.Background(), connect.NewRequest(msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}
