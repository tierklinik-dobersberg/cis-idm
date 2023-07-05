package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/mdp/qrterminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func getProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profile [sub-command]",
		Aliases: []string{"self", "self-service", "selfservice"},
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
		getAddAddressCommand(),
		getDeleteAddressCommand(),
		getUpdateAddressCommand(),
		getUpdateProfileCommand(),
		getEnrollTotpCommand(),
		getDisable2FACommand(),
		getGenerateRecoveryCodesCommand(),
		getSetAvatarCommand(),
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

func getAddAddressCommand() *cobra.Command {
	var msg idmv1.AddAddressRequest
	cmd := &cobra.Command{
		Use: "add-address",
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			_, err := cli.AddAddress(context.Background(), connect.NewRequest(&msg))
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

func getDeleteAddressCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete-address",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			_, err := cli.DeleteAddress(context.Background(), connect.NewRequest(&idmv1.DeleteAddressRequest{
				Id: args[0],
			}))

			if err != nil {
				logrus.Fatal(err.Error())
			}
		},
	}

	return cmd
}

func getUpdateAddressCommand() *cobra.Command {
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

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			_, err := cli.UpdateAddress(context.Background(), connect.NewRequest(&msg))
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

func getUpdateProfileCommand() *cobra.Command {
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

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)
			_, err := cli.UpdateProfile(context.Background(), connect.NewRequest(&msg))
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

func getEnrollTotpCommand() *cobra.Command {
	var (
		displayQR bool
	)

	cmd := &cobra.Command{
		Use: "enroll-totp",
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			res, err := cli.Enroll2FA(context.Background(), connect.NewRequest(&idmv1.Enroll2FARequest{
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

			_, err = cli.Enroll2FA(context.Background(), connect.NewRequest(&idmv1.Enroll2FARequest{
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

func getDisable2FACommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "disable-totp",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("\nPlease enter code: ")
			code, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				logrus.Fatal(err)
			}

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			_, err = cli.Remove2FA(context.Background(), connect.NewRequest(&idmv1.Remove2FARequest{
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

func getGenerateRecoveryCodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "generate-recovery-codes",
		Run: func(cmd *cobra.Command, args []string) {
			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			res, err := cli.GenerateRecoveryCodes(context.Background(), connect.NewRequest(&idmv1.GenerateRecoveryCodesRequest{}))
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

func getSetAvatarCommand() *cobra.Command {
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

			cli := idmv1connect.NewSelfServiceServiceClient(httpClient, baseURL)

			_, err := cli.UpdateProfile(context.Background(), connect.NewRequest(&idmv1.UpdateProfileRequest{
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
