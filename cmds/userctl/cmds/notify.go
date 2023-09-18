package cmds

import (
	"context"
	"encoding/json"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetSendNotificationCommand(root *cli.Root) *cobra.Command {
	var (
		targetRoleIDs   []string
		targetUserIDs   []string
		targetUserNames []string

		subject     string
		body        string
		bodyFile    string
		contextFile string
	)

	cmd := &cobra.Command{
		Use:     "send",
		Aliases: []string{"sms", "mail"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			targetUsers := make(map[string]struct{})

			for _, usrId := range targetUserIDs {
				targetUsers[usrId] = struct{}{}
			}

			for _, usrName := range targetUserNames {
				usr, err := root.Users().GetUser(ctx, connect.NewRequest(&idmv1.GetUserRequest{
					Search: &idmv1.GetUserRequest_Name{
						Name: usrName,
					},
					FieldMask: &fieldmaskpb.FieldMask{
						Paths: []string{"profile.user.id"},
					},
				}))
				if err != nil {
					logrus.Fatal(err)
				}

				targetUsers[usr.Msg.Profile.User.Id] = struct{}{}
			}

			if bodyFile != "" {
				if body != "" {
					logrus.Fatal("only one of --body or --body-file must be set")
				}
				content, err := os.ReadFile(bodyFile)
				if err != nil {
					logrus.Fatal(err)
				}

				body = string(content)
			}

			var (
				tmplCtx map[string]map[string]any
			)
			switch contextFile {
			case "":
			case "-":
				dec := json.NewDecoder(os.Stdin)
				if err := dec.Decode(&tmplCtx); err != nil {
					logrus.Fatal(err)
				}
			default:
				content, err := os.ReadFile(contextFile)
				if err != nil {
					logrus.Fatal(err)
				}

				if err := json.Unmarshal(content, &tmplCtx); err != nil {
					logrus.Fatal(err)
				}
			}

			req := &idmv1.SendNotificationRequest{}

			if tmplCtx != nil {
				req.PerUserTemplateContext = make(map[string]*structpb.Struct)

				var err error

				for key, putc := range tmplCtx {
					req.PerUserTemplateContext[key], err = structpb.NewStruct(putc)
					if err != nil {
						logrus.Fatal(err)
					}

				}
			}

			cli := idmv1connect.NewNotifyServiceClient(root.HttpClient, root.BaseURLS.Idm)

			switch cmd.CalledAs() {
			case "mail":
				req.Message = &idmv1.SendNotificationRequest_Email{
					Email: &idmv1.EMailMessage{
						Subject: subject,
						Body:    body,
					},
				}
			case "sms":
				req.Message = &idmv1.SendNotificationRequest_Sms{
					Sms: &idmv1.SMS{
						Body: body,
					},
				}
			default:
				logrus.Fatalf("please use the 'notify' command using the alias 'sms' or 'mail'")
			}

			targetUserSlice := make([]string, 0, len(targetUsers))
			for usr := range targetUsers {
				targetUserSlice = append(targetUserSlice, usr)
			}

			req.TargetRoles = targetRoleIDs
			req.TargetUsers = targetUserSlice

			res, err := cli.SendNotification(ctx, connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err)
			}

			root.Print(res.Msg)
		},
	}

	flags := cmd.Flags()
	{
		flags.StringSliceVar(&targetRoleIDs, "to-role-ids", nil, "A list of role IDs that should receive the notification")
		flags.StringSliceVar(&targetUserIDs, "to-user-ids", nil, "A list of user IDs that should receive the notification")
		flags.StringSliceVar(&targetUserNames, "to-user", nil, "A list of usernames that should receive the notificatoin")

		flags.StringVar(&subject, "subject", "", "The subject if sending a mail")
		flags.StringVar(&body, "body", "", "The body of the email/sms")
		flags.StringVar(&bodyFile, "body-file", "", "Path to the file that contains the body for the SMS/mail")
	}

	return cmd
}
