package cmds

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetUsersCommand(root *cli.Root) *cobra.Command {
	var (
		fm            []string
		excludeFields bool
	)

	cmd := &cobra.Command{
		Use: "users",
		Run: func(cmd *cobra.Command, args []string) {
			req := &idmv1.ListUsersRequest{}

			if len(fm) > 0 {
				req.FieldMask = &fieldmaskpb.FieldMask{}
				for _, name := range fm {
					req.FieldMask.Paths = append(req.FieldMask.Paths, fmt.Sprintf("users.%s", name))
				}
			}
			req.ExcludeFields = excludeFields

			users, err := root.Users().ListUsers(context.Background(), connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err.Error())
			}

			root.Print(users.Msg.Users)
		},
	}

	cmd.Flags().StringSliceVar(&fm, "fields", nil, "A list of fields to include")
	cmd.Flags().BoolVar(&excludeFields, "exclude-fields", false, "Use --fields as an exclude list rather than an include list")

	cmd.AddCommand(
		GetDeleteUserCommand(root),
		GetGenerateRegistrationTokenCommand(root),
		GetInviteUserCommand(root),
		GetUpdateUserCommand(root),
		GetCreateUserCommand(root),
	)

	return cmd
}

func GetDeleteUserCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, err := root.Users().DeleteUser(context.Background(), connect.NewRequest(&idmv1.DeleteUserRequest{
				Id: args[0],
			}))

			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func GetRegisterUserCommand(root *cli.Root) *cobra.Command {
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

			msg := &idmv1.RegisterUserRequest{
				Username:          args[0],
				Password:          password,
				RegistrationToken: regToken,
			}

			res, err := root.Auth().RegisterUser(context.Background(), connect.NewRequest(msg))
			if err != nil {
				logrus.Fatal(err)
			}

			if res.Msg.GetAccessToken() == nil {
				logrus.Fatal("unexpected server response received.")
			}

			root.Print(res.Msg)
		},
	}

	cmd.Flags().StringVar(&password, "password", "", "The password for the new user account.")
	cmd.Flags().StringVar(&regToken, "registration-token", "", "The registration token to authenticate the request.")

	return cmd
}

func GetGenerateRegistrationTokenCommand(root *cli.Root) *cobra.Command {
	var (
		ttl      time.Duration
		maxUsage int
	)

	cmd := &cobra.Command{
		Use: "generate-registration-token",
		Run: func(cmd *cobra.Command, args []string) {
			msg := &idmv1.GenerateRegistrationTokenRequest{}

			if ttl > 0 {
				msg.Ttl = durationpb.New(ttl)
			}

			if maxUsage > 0 {
				msg.MaxCount = uint64(maxUsage)
			}

			req := connect.NewRequest(msg)
			res, err := root.Auth().GenerateRegistrationToken(context.Background(), req)
			if err != nil {
				logrus.Fatal(err)
			}

			root.Print(res.Msg)
		},
	}

	f := cmd.Flags()
	{
		f.DurationVar(&ttl, "ttl", 0, "The time-to-live for the access token")
		f.IntVar(&maxUsage, "max-usage", 1, "How often the token can be used")
	}

	return cmd
}

func GetInviteUserCommand(root *cli.Root) *cobra.Command {
	var roles []string
	cmd := &cobra.Command{
		Use:  "invite",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			req := &idmv1.InviteUserRequest{
				InitialRoles: roles,
			}

			for _, arg := range args {
				parts := strings.Split(arg, ":")
				if len(parts) > 1 {
					req.Invite = append(req.Invite, &idmv1.UserInvite{
						Name:  parts[1],
						Email: parts[0],
					})
				} else {
					req.Invite = append(req.Invite, &idmv1.UserInvite{
						Name:  parts[0],
						Email: parts[0],
					})
				}
			}

			_, err := root.Users().InviteUser(context.Background(), connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&roles, "roles", nil, "Automatically assign the given roles upon user registration")

	return cmd
}

func GetCreateUserCommand(root *cli.Root) *cobra.Command {
	var (
		user      = new(idmv1.User)
		emails    []string
		phone     []string
		password  string
		bcrypt    bool
		roleIds   []string
		roleNames []string
	)

	cmd := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			req := &idmv1.CreateUserRequest{
				Profile: &idmv1.Profile{
					User: user,
				},
				Password:         password,
				PasswordIsBcrypt: bcrypt,
			}

			for _, m := range emails {
				req.Profile.EmailAddresses = append(req.Profile.EmailAddresses, &idmv1.EMail{
					Address: m,
				})
			}

			for idx, p := range phone {
				req.Profile.PhoneNumbers = append(req.Profile.PhoneNumbers, &idmv1.PhoneNumber{
					Number:   p,
					Verified: true,
					Primary:  idx == 0,
				})
			}

			for _, id := range roleIds {
				req.Profile.Roles = append(req.Profile.Roles, &idmv1.Role{
					Id: id,
				})
			}

			for _, name := range roleNames {
				req.Profile.Roles = append(req.Profile.Roles, &idmv1.Role{
					Name: name,
				})
			}

			res, err := root.Users().CreateUser(context.Background(), connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err)
			}

			root.Print(res.Msg)
		},
	}

	f := cmd.Flags()
	{
		f.StringVar(&user.FirstName, "first-name", "", "")
		f.StringVar(&user.LastName, "last-name", "", "")
		f.StringVar(&user.Username, "name", "", "")
		f.StringVar(&user.DisplayName, "display-name", "", "")
		f.StringVar(&password, "password", "", "")
		f.BoolVar(&bcrypt, "bcrypt", false, "")
		f.StringSliceVar(&emails, "email", nil, "")
		f.StringSliceVar(&phone, "phone", nil, "")
		f.StringSliceVar(&roleIds, "role-id", nil, "")
		f.StringSliceVar(&roleNames, "role", nil, "")
	}

	return cmd
}

func GetUpdateUserCommand(root *cli.Root) *cobra.Command {
	var (
		msg       idmv1.UpdateUserRequest
		extraData []string
	)

	cmd := &cobra.Command{
		Use:  "update",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg.FieldMask = &fieldmaskpb.FieldMask{}
			msg.Id = args[0]

			flagsToFieldmask := [][2]string{
				{"username", "username"},
				{"display-name", "display_name"},
				{"first-name", "first_name"},
				{"last-name", "last_name"},
				{"avatar", "avatar"},
				{"birthday", "birthday"},
				{"extra", "extra"},
			}

			for _, flag := range flagsToFieldmask {
				if cmd.Flag(flag[0]).Changed {
					msg.FieldMask.Paths = append(msg.FieldMask.Paths, flag[1])
				}
			}

			if len(extraData) > 0 {
				extra := make(map[string]any)

				for _, value := range extraData {
					parts := strings.SplitN(value, ":", 3)
					if len(parts) == 2 {
						parts = []string{parts[0], "string", parts[1]}
					}

					if len(parts) != 3 {
						logrus.Fatalf("invalid value for --extra, expected format <key>:<type>:<value> or <key>:<value>: %q", parts)
					}

					var parsed any
					switch parts[1] {
					case "string", "":
						if parts[2] != "null" {
							parsed = parts[2]
						}
					case "int":
						intVal, err := strconv.ParseInt(parts[2], 10, 0)
						if err != nil {
							logrus.Fatal(err)
						}
						parsed = intVal
					case "float":
						floatVal, err := strconv.ParseFloat(parts[2], 0)
						if err != nil {
							logrus.Fatal(err)
						}
						parsed = floatVal
					case "bool":
						boolVal, err := strconv.ParseBool(parts[2])
						if err != nil {
							logrus.Fatal(err)
						}
						parsed = boolVal
					default:
						logrus.Fatalf("unsupported type for extra data %q", parts[1])
					}

					extra[parts[0]] = parsed
				}

				if len(extra) > 0 {
					var err error
					msg.Extra, err = structpb.NewStruct(extra)
					if err != nil {
						logrus.Fatal(err)
					}
				}
			}

			_, err := root.Users().UpdateUser(context.Background(), connect.NewRequest(&msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVar(&msg.Username, "username", "", "The new user name")
		flags.StringVar(&msg.DisplayName, "display-name", "", "The new display name")
		flags.StringVar(&msg.FirstName, "first-name", "", "The new first name")
		flags.StringVar(&msg.LastName, "last-name", "", "The new last name")
		flags.StringVar(&msg.Avatar, "avatar", "", "The new avatar value")
		flags.StringVar(&msg.Birthday, "birthday", "", "The birthday of the user")
		flags.StringSliceVar(&extraData, "extra", nil, "Format is key:type:value")
	}

	return cmd
}
