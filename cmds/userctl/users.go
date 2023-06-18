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
