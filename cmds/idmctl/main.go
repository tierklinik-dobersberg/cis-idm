package main

import (
	"errors"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
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
				var diag hcl.Diagnostics
				if errors.As(err, &diag) {
					merr := new(multierror.Error)
					for _, d := range diag {
						merr.Errors = append(merr.Errors, d)
					}

					logrus.Fatalf(merr.Error())
				}

				logrus.Fatalf(err.Error())
			}

			root.Print(c)
		},
	})

	if err := root.Execute(); err != nil {
		logrus.Fatalf("failed to run: %s", err)
	}
}
