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
	"google.golang.org/protobuf/types/known/structpb"
)

func GetSendNotificationCommand(root *cli.Root) *cobra.Command {
	var (
		targetUser []string

		subject     string
		body        string
		bodyFile    string
		contextFile string
		webpushRaw  bool
	)

	cmd := &cobra.Command{
		Use:     "send",
		Aliases: []string{"sms", "mail", "web-push"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			targetUsers := root.MustResolveUserIds(targetUser)
			logrus.Infof("target-users: %v", targetUsers)

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

			cli := idmv1connect.NewNotifyServiceClient(root.HttpClient, root.Config().BaseURLS.Idm)

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

			case "web-push":
				var blob []byte

				if !webpushRaw {
					// TODO(ppacher): this is just testing code,
					// remove and make that an actual usable thingy...
					notify := map[string]any{
						"notification": map[string]any{
							"title": "CIS",
							"body":  body,
							"actions": []any{
								map[string]any{
									"action": "foo",
									"title":  "Foo",
								},
								map[string]any{
									"action": "bar",
									"title":  "Bar",
								},
							},
							"data": map[string]any{
								"onActionClick": map[string]any{
									"default": map[string]any{"operation": "openWindow"},
									"foo":     map[string]any{"operation": "openWindow", "url": "https://account.dobersberg.vet"},
								},
							},
						},
					}

					blob, _ = json.Marshal(notify)
				} else {
					blob = []byte(body)
				}

				req.Message = &idmv1.SendNotificationRequest_Webpush{
					Webpush: &idmv1.WebPushNotification{
						Kind: &idmv1.WebPushNotification_Template{
							Template: string(blob),
						},
					},
				}

			default:
				logrus.Fatalf("please use the 'notify' command using the alias 'sms' or 'mail'")
			}

			req.TargetUsers = targetUsers

			res, err := cli.SendNotification(ctx, connect.NewRequest(req))
			if err != nil {
				logrus.Fatal(err)
			}

			root.Print(res.Msg)
		},
	}

	flags := cmd.Flags()
	{
		flags.StringSliceVar(&targetUser, "to-user", nil, "A list of usernames or ids that should receive the notificatoin")

		flags.StringVar(&subject, "subject", "", "The subject if sending a mail")
		flags.StringVar(&body, "body", "", "The body of the email/sms")
		flags.StringVar(&bodyFile, "body-file", "", "Path to the file that contains the body for the SMS/mail")
		flags.BoolVar(&webpushRaw, "web-push-raw", false, "Do not convert body into a notification object")
	}

	return cmd
}
