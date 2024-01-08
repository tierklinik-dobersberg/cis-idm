package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"github.com/tierklinik-dobersberg/cis-idm/cmds/idmctl/cmds"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
)

func main() {
	root := cli.New("idmctl")

	cmds.PrepareRootCommand(root)

	root.AddCommand(&cobra.Command{
		Use:  "validate-config",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := config.LoadFile(args[0])
			if err != nil {
				logrus.Fatalf(err.Error())
			}

			root.Print(c)
		},
	})

	if err := root.Execute(); err != nil {
		logrus.Fatalf("failed to run: %s", err)
	}
}
