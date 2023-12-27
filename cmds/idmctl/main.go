package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
	"github.com/tierklinik-dobersberg/cis-idm/cmds/idmctl/cmds"
)

func main() {
	root := cli.New("idmctl")

	cmds.PrepareRootCommand(root)

	if err := root.Execute(); err != nil {
		logrus.Fatalf("failed to run: %s", err)
	}
}