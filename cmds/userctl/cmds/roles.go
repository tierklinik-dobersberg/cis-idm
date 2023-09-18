package cmds

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
)

func GetRoleCommand(root *cli.Root) *cobra.Command {
	var (
		byName bool
	)

	cmd := &cobra.Command{
		Use:     "roles",
		Aliases: []string{"role"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				res, err := root.Roles().ListRoles(context.Background(), connect.NewRequest(new(idmv1.ListRolesRequest)))
				if err != nil {
					logrus.Fatal(err)
				}

				root.Print(res.Msg.Roles)

				return
			}

			req := &idmv1.GetRoleRequest{}
			if byName {
				req.Search = &idmv1.GetRoleRequest_Name{
					Name: args[0],
				}
			} else {
				req.Search = &idmv1.GetRoleRequest_Id{
					Id: args[0],
				}
			}

			res, err := root.Roles().GetRole(context.Background(), connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err)
			}

			root.Print(res.Msg)
		},
	}

	cmd.Flags().BoolVar(&byName, "name", false, "Search role by name")

	cmd.AddCommand(
		GetCreateRoleCommand(root),
		GetDeleteRoleCommand(root),
		GetAssignRoleCommand(root),
		GetUnassignRoleCommand(root),
	)

	return cmd
}

func GetCreateRoleCommand(root *cli.Root) *cobra.Command {
	var (
		name            string
		description     string
		deleteProtected bool
	)

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Run: func(cmd *cobra.Command, args []string) {
			res, err := root.Roles().CreateRole(context.Background(), connect.NewRequest(&idmv1.CreateRoleRequest{
				Name:             name,
				Description:      description,
				DeleteProtection: deleteProtected,
			}))
			if err != nil {
				logrus.Fatal(err)
			}

			root.Print(res.Msg.Role)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "The name for the new role")
	cmd.Flags().StringVar(&description, "description", "", "The description for the new role")
	cmd.Flags().BoolVar(&deleteProtected, "delete-protected", false, "Whether or not the new role should be delete protected")

	cmd.MarkFlagRequired("name")

	return cmd
}

func GetDeleteRoleCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, err := root.Roles().DeleteRole(context.Background(), connect.NewRequest(&idmv1.DeleteRoleRequest{
				RoleId: args[0],
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	return cmd
}

func GetAssignRoleCommand(root *cli.Root) *cobra.Command {
	var (
		userIDs    []string
		roleByName bool
	)

	cmd := &cobra.Command{
		Use:  "assign",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var roleID string
			if roleByName {
				res, err := root.Roles().GetRole(context.Background(), connect.NewRequest(&idmv1.GetRoleRequest{
					Search: &idmv1.GetRoleRequest_Name{
						Name: args[0],
					},
				}))
				if err != nil {
					logrus.Fatal(err)
				}

				roleID = res.Msg.Role.Id
			} else {
				roleID = args[0]
			}

			_, err := root.Roles().AssignRoleToUser(context.Background(), connect.NewRequest(&idmv1.AssignRoleToUserRequest{
				RoleId: roleID,
				UserId: userIDs,
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&userIDs, "user", nil, "A list of users ids to assign the role to")
	cmd.Flags().BoolVar(&roleByName, "role-name", false, "Search for the role by name first")

	return cmd
}

func GetUnassignRoleCommand(root *cli.Root) *cobra.Command {
	var (
		userIDs    []string
		roleByName bool
	)

	cmd := &cobra.Command{
		Use:  "unassign",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var roleID string
			if roleByName {
				res, err := root.Roles().GetRole(context.Background(), connect.NewRequest(&idmv1.GetRoleRequest{
					Search: &idmv1.GetRoleRequest_Name{
						Name: args[0],
					},
				}))
				if err != nil {
					logrus.Fatal(err)
				}

				roleID = res.Msg.Role.Id
			} else {
				roleID = args[0]
			}

			_, err := root.Roles().UnassignRoleFromUser(context.Background(), connect.NewRequest(&idmv1.UnassignRoleFromUserRequest{
				RoleId: roleID,
				UserId: userIDs,
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&userIDs, "user", nil, "A list of users ids to assign the role to")
	cmd.Flags().BoolVar(&roleByName, "role-name", false, "Search for the role by name first")

	return cmd
}
