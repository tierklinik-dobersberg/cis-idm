package cmds

import (
	"fmt"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/mdp/qrterminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetProfileCommand(root *cli.Root) *cobra.Command {
	var (
		readMaskPaths []string
		excludeFields bool
	)

	cmd := &cobra.Command{
		Use:     "profile [sub-command]",
		Aliases: []string{"self", "self-service", "selfservice"},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := &idmv1.IntrospectRequest{}

			if len(readMaskPaths) > 0 {
				req.ReadMask = &fieldmaskpb.FieldMask{
					Paths: readMaskPaths,
				}

				req.ExcludeFields = excludeFields
			} else if excludeFields {
				logrus.Fatalf("--exclude-fields can only be used if --fields is specified")
			}

			res, err := root.Auth().Introspect(root.Context(), connect.NewRequest(req))
			if err != nil {
				return err
			}

			if profile := res.Msg.Profile; profile != nil {
				root.Print(profile)
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&readMaskPaths, "fields", nil, "Include/Exclude specified fields")
	cmd.Flags().BoolVar(&excludeFields, "exclude-fields", false, "Use --fields for exclusion rather than inclusion")

	cmd.AddCommand(
		GetChangePasswordCommand(root),
		GetAddEmailCommand(root),
		GetDeleteEmailCommand(root),
		GetAddAddressCommand(root),
		GetDeleteAddressCommand(root),
		GetUpdateAddressCommand(root),
		GetUpdateProfileCommand(root),
		GetEnrollTotpCommand(root),
		GetDisable2FACommand(root),
		GetGenerateRecoveryCodesCommand(root),
		GetSetAvatarCommand(root),
		GetAPITokenCommand(root),
	)

	return cmd
}

func GetAPITokenCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use: "api-token",
	}

	cmd.AddCommand(
		GetGenerateAPITokenCommand(root),
		GetListAPITokensCommand(root),
		GetRevokeAPITokenCommand(root),
	)

	return cmd
}

func GetGenerateAPITokenCommand(root *cli.Root) *cobra.Command {
	var (
		expiresAt string
		roles     []string
	)

	cmd := &cobra.Command{
		Use:  "generate [name]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			roleIds := root.MustResolveRoleIds(roles)

			req := &idmv1.GenerateAPITokenRequest{
				Roles:       roleIds,
				Description: args[0],
			}

			if expiresAt != "" {
				e, err := time.Parse(time.RFC3339, expiresAt)
				if err != nil {
					logrus.Fatalf("invalid value for --expires: %s", err)
				}

				req.Expires = timestamppb.New(e)
			}

			res, err := root.SelfService().GenerateAPIToken(root.Context(), connect.NewRequest(req))

			if err != nil {
				logrus.Fatal(err.Error())
			}

			root.Print(res)
		},
	}

	cmd.Flags().StringVar(&expiresAt, "expires", "", "A Timestamp in RFC3339 at which the token should expire")
	cmd.Flags().StringSliceVar(&roles, "role", nil, "A list of roles to assign")

	return cmd
}

func GetListAPITokensCommand(root *cli.Root) *cobra.Command {
	return &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			res, err := root.SelfService().ListAPITokens(
				root.Context(),
				connect.NewRequest(new(idmv1.ListAPITokensRequest)),
			)
			if err != nil {
				logrus.Fatalf(err.Error())
			}

			root.Print(res.Msg)
		},
	}
}

func GetRevokeAPITokenCommand(root *cli.Root) *cobra.Command {
	return &cobra.Command{
		Use:     "delete [token-id]",
		Aliases: []string{"revoke"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res, err := root.SelfService().RemoveAPIToken(
				root.Context(),
				connect.NewRequest(&idmv1.RemoveAPITokenRequest{
					Id: args[0],
				}),
			)

			if err != nil {
				logrus.Fatalf(err.Error())
			}

			root.Print(res.Msg)
		},
	}
}

func GetChangePasswordCommand(root *cli.Root) *cobra.Command {
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

			_, err := root.SelfService().ChangePassword(root.Context(), connect.NewRequest(&req))
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

func GetAddEmailCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-mail",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := &idmv1.AddEmailAddressRequest{
				Email: args[0],
			}

			_, err := root.SelfService().AddEmailAddress(root.Context(), connect.NewRequest(msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func GetDeleteEmailCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete-mail",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg := &idmv1.DeleteEmailAddressRequest{
				Id: args[0],
			}

			_, err := root.SelfService().DeleteEmailAddress(root.Context(), connect.NewRequest(msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func GetAddAddressCommand(root *cli.Root) *cobra.Command {
	var msg idmv1.AddAddressRequest
	cmd := &cobra.Command{
		Use: "add-address",
		Run: func(cmd *cobra.Command, args []string) {

			_, err := root.SelfService().AddAddress(root.Context(), connect.NewRequest(&msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVar(&msg.CityCode, "city-code", "", "The city code of the new address")
		flags.StringVar(&msg.CityName, "city", "", "The name of the city of the new address")
		flags.StringVar(&msg.Street, "street", "", "The street of the new address")
		flags.StringVar(&msg.Extra, "extra", "", "An additional specifier for the new address")
	}

	return cmd
}

func GetDeleteAddressCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete-address",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, err := root.SelfService().DeleteAddress(root.Context(), connect.NewRequest(&idmv1.DeleteAddressRequest{
				Id: args[0],
			}))

			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func GetUpdateAddressCommand(root *cli.Root) *cobra.Command {
	var msg idmv1.UpdateAddressRequest
	cmd := &cobra.Command{
		Use:  "update-address",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			msg.Id = args[0]
			msg.FieldMask = &fieldmaskpb.FieldMask{}

			if cmd.Flag("city-code").Changed {
				msg.FieldMask.Paths = append(msg.FieldMask.Paths, "city_code")
			}
			if cmd.Flag("city").Changed {
				msg.FieldMask.Paths = append(msg.FieldMask.Paths, "city_name")
			}
			if cmd.Flag("street").Changed {
				msg.FieldMask.Paths = append(msg.FieldMask.Paths, "street")
			}
			if cmd.Flag("extra").Changed {
				msg.FieldMask.Paths = append(msg.FieldMask.Paths, "extra")
			}

			_, err := root.SelfService().UpdateAddress(root.Context(), connect.NewRequest(&msg))
			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	flags := cmd.Flags()
	{
		flags.StringVar(&msg.CityCode, "city-code", "", "The city code of the new address")
		flags.StringVar(&msg.CityName, "city", "", "The name of the city of the new address")
		flags.StringVar(&msg.Street, "street", "", "The street of the new address")
		flags.StringVar(&msg.Extra, "extra", "", "An additional specifier for the new address")
	}

	return cmd
}

func GetUpdateProfileCommand(root *cli.Root) *cobra.Command {
	var msg idmv1.UpdateProfileRequest
	cmd := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			msg.FieldMask = &fieldmaskpb.FieldMask{}

			flagsToFieldmask := [][2]string{
				{"username", "username"},
				{"display-name", "display_name"},
				{"first-name", "first_name"},
				{"last-name", "last_name"},
				{"avatar", "avatar"},
				{"birthday", "birthday"},
			}

			for _, flag := range flagsToFieldmask {
				if cmd.Flag(flag[0]).Changed {
					msg.FieldMask.Paths = append(msg.FieldMask.Paths, flag[1])
				}
			}

			_, err := root.SelfService().UpdateProfile(root.Context(), connect.NewRequest(&msg))
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
	}

	return cmd
}

func GetEnrollTotpCommand(root *cli.Root) *cobra.Command {
	var (
		displayQR bool
	)

	cmd := &cobra.Command{
		Use: "enroll-totp",
		Run: func(cmd *cobra.Command, args []string) {

			res, err := root.SelfService().Enroll2FA(root.Context(), connect.NewRequest(&idmv1.Enroll2FARequest{
				Kind: &idmv1.Enroll2FARequest_TotpStep1{},
			}))
			if err != nil {
				logrus.Fatal(err)
			}

			if displayQR {
				qrterminal.Generate(res.Msg.GetTotpStep1().Url, qrterminal.L, os.Stdout)
			} else {
				fmt.Printf("Secret: %s\nUrl: %s\n", res.Msg.GetTotpStep1().Secret, res.Msg.GetTotpStep1().Url)
			}

			fmt.Print("\nPlease enter code: ")
			code, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				logrus.Fatal(err)
			}

			_, err = root.SelfService().Enroll2FA(root.Context(), connect.NewRequest(&idmv1.Enroll2FARequest{
				Kind: &idmv1.Enroll2FARequest_TotpStep2{
					TotpStep2: &idmv1.EnrollTOTPRequestStep2{
						VerifyCode: string(code),
						Secret:     res.Msg.GetTotpStep1().Secret,
						SecretHmac: res.Msg.GetTotpStep1().GetSecretHmac(),
					},
				},
			}))
			if err != nil {
				logrus.Fatal(err)
			}

		},
	}

	cmd.Flags().BoolVar(&displayQR, "qr", true, "Display a QR text")

	return cmd
}

func GetDisable2FACommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use: "disable-totp",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("\nPlease enter code: ")
			code, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				logrus.Fatal(err)
			}

			_, err = root.SelfService().Remove2FA(root.Context(), connect.NewRequest(&idmv1.Remove2FARequest{
				Kind: &idmv1.Remove2FARequest_TotpCode{
					TotpCode: string(code),
				},
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	return cmd
}

func GetGenerateRecoveryCodesCommand(root *cli.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use: "generate-recovery-codes",
		Run: func(cmd *cobra.Command, args []string) {

			res, err := root.SelfService().GenerateRecoveryCodes(root.Context(), connect.NewRequest(&idmv1.GenerateRecoveryCodesRequest{}))
			if err != nil {
				logrus.Fatal(err)
			}

			l := len(res.Msg.RecoveryCodes)
			fmt.Printf("Recovery Codes (%d): \n", l)

			for i := 0; i < (l / 4); i++ {
				for j := 0; j < 4 && (i*4+j) < l; j++ {
					idx := i*4 + j
					fmt.Printf(" %s ", res.Msg.RecoveryCodes[idx])
				}
				fmt.Println()
			}
		},
	}

	return cmd
}

func GetSetAvatarCommand(root *cli.Root) *cobra.Command {
	var pathIsURL bool

	cmd := &cobra.Command{
		Use:  "set-avatar",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			if !pathIsURL {
				content, err := os.ReadFile(args[0])
				if err != nil {
					logrus.Fatal(err)
				}

				url = dataurl.EncodeBytes(content)
			}

			_, err := root.SelfService().UpdateProfile(root.Context(), connect.NewRequest(&idmv1.UpdateProfileRequest{
				Avatar: url,
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"avatar"},
				},
			}))
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVar(&pathIsURL, "url", false, "The specified path is a remote URL")

	return cmd
}
