package cmds

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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
		GetImportRolesCommand(root),
		GetCreateRoleCommand(root),
		GetUpdateRoleCommand(root),
		GetDeleteRoleCommand(root),
		GetAssignRoleCommand(root),
		GetUnassignRoleCommand(root),
		GetResolveRolePermissions(root),
	)

	return cmd
}

func GetResolveRolePermissions(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get-permissions [role-id/name]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			role, err := root.ResolveRole(args[0])
			if err != nil {
				logrus.Fatalf(err.Error())
			}

			res, err := root.Roles().ResolveRolePermissions(root.Context(), connect.NewRequest(&idmv1.ResolveRolePermissionsRequest{
				RoleId: role.Id,
			}))
			if err != nil {
				logrus.Fatalf(err.Error())
			}

			root.Print(res.Msg)
		},
	}

	return cmd
}

func GetImportRolesCommand(root *cli.Root) *cobra.Command {
	var (
		ignoreExisting bool
	)

	cmd := &cobra.Command{
		Use:  "import",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := os.Stdin
			if args[0] != "-" {
				var err error
				file, err = os.Open(args[0])
				if err != nil {
					logrus.Fatalf("failed to open %q: %s", args[0], err)
				}

				defer file.Close()
			}

			content, err := io.ReadAll(file)
			if err != nil {
				logrus.Fatalf("failed to read from %q: %s", args[0], err)
			}

			var listResponse []*idmv1.Role
			if err := json.Unmarshal(content, &listResponse); err != nil {
				logrus.Fatalf("failed to parse ListRolesResponse: %s", err)
			}

			cli := root.Roles()
			ctx := root.Context()

			for _, role := range listResponse {
				if ignoreExisting {
					// check if a role with the same ID or name exists and skip it otherwise
					existing, err := root.ResolveRole(role.Id)
					if err != nil {
						logrus.Fatalf("failed to check if role id=%q name=%q exists: %s", role.Id, role.Name, err)
					}

					if existing != nil {
						logrus.Infof("skipping import of existing role id=%q", role.Id)
						continue
					}

					existing, err = root.ResolveRole(role.Name)
					if err != nil {
						logrus.Fatalf("failed to check if role id=%q name=%q exists: %s", role.Id, role.Name, err)
					}

					if existing != nil {
						logrus.Infof("skipping import of existing role name=%q", role.Name)
						continue
					}
				}

				_, err := cli.CreateRole(ctx, connect.NewRequest(&idmv1.CreateRoleRequest{
					Id:               role.Id,
					Name:             role.Name,
					Description:      role.Description,
					DeleteProtection: role.DeleteProtected,
				}))

				if err != nil {
					logrus.Fatalf("failed to create role id=%q name=%q: %s", role.Id, role.Name, err)
				}
			}
		},
	}

	cmd.Flags().BoolVar(&ignoreExisting, "ignore-existing", false, "Do not try to import roles where the ID or Name is already used")

	return cmd
}

func GetCreateRoleCommand(root *cli.Root) *cobra.Command {
	var (
		id              string
		name            string
		description     string
		deleteProtected bool
	)

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Run: func(cmd *cobra.Command, args []string) {
			res, err := root.Roles().CreateRole(context.Background(), connect.NewRequest(&idmv1.CreateRoleRequest{
				Id:               id,
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

	cmd.Flags().StringVar(&id, "id", "", "The ID for the new role. If unset, a random UUID will be generated")
	cmd.Flags().StringVar(&name, "name", "", "The name for the new role")
	cmd.Flags().StringVar(&description, "description", "", "The description for the new role")
	cmd.Flags().BoolVar(&deleteProtected, "delete-protected", false, "Whether or not the new role should be delete protected")

	cmd.MarkFlagRequired("name")

	return cmd
}

func GetUpdateRoleCommand(root *cli.Root) *cobra.Command {
	var (
		name            string
		description     string
		deleteProtected bool
	)

	cmd := &cobra.Command{
		Use:  "update",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			role, err := root.ResolveRole(args[0])
			if err != nil {
				logrus.Fatalf("failed to resolve role: %s", err)
			}

			fm := &fieldmaskpb.FieldMask{}

			cases := [][]string{
				{"name", "name"},
				{"description", "description"},
				{"deleted_protection", "delete-protected"},
			}

			for _, s := range cases {
				if cmd.Flag(s[1]).Changed {
					fm.Paths = append(fm.Paths, s[0])
				}
			}

			if len(fm.Paths) == 0 {
				logrus.Fatalf("nothing changed")
			}

			res, err := root.Roles().UpdateRole(context.Background(), connect.NewRequest(&idmv1.UpdateRoleRequest{
				RoleId:           role.Id,
				Name:             name,
				Description:      description,
				DeleteProtection: deleteProtected,
				FieldMask:        fm,
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
		users []string
	)

	cmd := &cobra.Command{
		Use:  "assign",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			role, err := root.ResolveRole(args[0])
			if err != nil {
				logrus.Fatalf("failed to resolve role: %s", err)
			}

			users = root.MustResolveUserIds(users)

			_, err = root.Roles().AssignRoleToUser(context.Background(), connect.NewRequest(&idmv1.AssignRoleToUserRequest{
				RoleId: role.Id,
				UserId: users,
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&users, "to", nil, "A list of users ids to assign the role to")

	cmd.MarkFlagRequired("user")

	return cmd
}

func GetUnassignRoleCommand(root *cli.Root) *cobra.Command {
	var (
		users []string
	)

	cmd := &cobra.Command{
		Use:  "unassign",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			role, err := root.ResolveRole(args[0])
			if err != nil {
				logrus.Fatalf("failed to resolve role: %s", err)
			}

			users = root.MustResolveUserIds(users)

			_, err = root.Roles().UnassignRoleFromUser(context.Background(), connect.NewRequest(&idmv1.UnassignRoleFromUserRequest{
				RoleId: role.Id,
				UserId: users,
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&users, "from", nil, "A list of users ids to unassign the role")

	cmd.MarkFlagRequired("from")

	return cmd
}
