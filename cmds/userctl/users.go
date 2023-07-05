package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func getUsersCommand() *cobra.Command {
	var fm []string

	cmd := &cobra.Command{
		Use: "users",
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewUserServiceClient(httpClient, baseURL)

			req := &idmv1.ListUsersRequest{}

			if len(fm) > 0 {
				req.FieldMask = &fieldmaskpb.FieldMask{}
				for _, name := range fm {
					req.FieldMask.Paths = append(req.FieldMask.Paths, fmt.Sprintf("users.%s", name))
				}
			}
			users, err := cli.ListUsers(context.Background(), connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err.Error())
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")

			if err := enc.Encode(users.Msg.Users); err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringSliceVar(&fm, "fields", nil, "A list of fields to include")

	cmd.AddCommand(
		getDeleteUserCommand(),
		getGenerateRegistrationTokenCommand(),
	)

	return cmd
}

func getDeleteUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewUserServiceClient(httpClient, baseURL)

			_, err := cli.DeleteUser(context.Background(), connect.NewRequest(&idmv1.DeleteUserRequest{
				Id: args[0],
			}))

			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func getRegisterUserCommand() *cobra.Command {
	var (
		password string
		regToken string
	)

	cmd := &cobra.Command{
		Use:  "register",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if password == "" {
				fmt.Print("Please enter password: ")
				pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					logrus.Fatal(err)
				}

				fmt.Print("Please repeat password: ")
				pwd2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					logrus.Fatal(err)
				}

				if string(pwd) != string(pwd2) {
					logrus.Fatal("passwords do not match")
				}

				password = string(pwd)
			}

			cli := idmv1connect.NewAuthServiceClient(httpClient, baseURL)

			msg := &idmv1.RegisterUserRequest{
				Username:          args[0],
				Password:          password,
				RegistrationToken: regToken,
			}

			res, err := cli.RegisterUser(context.Background(), connect.NewRequest(msg))
			if err != nil {
				logrus.Fatal(err)
			}

			if res.Msg.GetAccessToken() == nil {
				logrus.Fatal("unexpected server response received.")
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")

			if err := enc.Encode(res.Msg.GetAccessToken()); err != nil {
				logrus.Fatal(err)
			}

		},
	}

	cmd.Flags().StringVar(&password, "password", "", "The password for the new user account.")
	cmd.Flags().StringVar(&regToken, "registration-token", "", "The registration token to authenticate the request.")

	return cmd
}

func getGenerateRegistrationTokenCommand() *cobra.Command {
	var (
		ttl      time.Duration
		maxUsage int
	)

	cmd := &cobra.Command{
		Use: "generate-registration-token",
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewAuthServiceClient(httpClient, baseURL)

			msg := &idmv1.GenerateRegistrationTokenRequest{}

			if ttl > 0 {
				msg.Ttl = durationpb.New(ttl)
			}

			if maxUsage > 0 {
				msg.MaxCount = uint64(maxUsage)
			}

			req := connect.NewRequest(msg)
			res, err := cli.GenerateRegistrationToken(context.Background(), req)
			if err != nil {
				logrus.Fatal(err)
			}

			fmt.Printf("Registration token: %s\n", res.Msg.Token)
		},
	}

	f := cmd.Flags()
	{
		f.DurationVar(&ttl, "ttl", 0, "The time-to-live for the access token")
		f.IntVar(&maxUsage, "max-usage", 1, "How often the token can be used")
	}

	return cmd
}
