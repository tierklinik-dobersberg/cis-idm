package cmds

import (
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
)

func PrepareRootCommand(root *cli.Root) {
	root.AddCommand(
		GetLoginCommand(root),
		GetProfileCommand(root),
		GetUsersCommand(root),
		GetRegisterUserCommand(root),
		GetRoleCommand(root),
		GetSendNotificationCommand(root),
	)
}
