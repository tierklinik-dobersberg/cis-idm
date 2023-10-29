package cmds

import (
	"errors"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/cli"
)

// resolveRole tries to load a IDM role either by ID or by Name. If the role does
// not exist a nil result and a nil error is returned. Otherwise, either the role or
// any other error is returned.
func resolveRole(root *cli.Root, roleNameOrId string) (*idmv1.Role, error) {
	cli := root.Roles()

	role, err := cli.GetRole(root.Context(), connect.NewRequest(&idmv1.GetRoleRequest{
		Search: &idmv1.GetRoleRequest_Id{
			Id: roleNameOrId,
		},
	}))

	var cerr *connect.Error
	if err == nil {
		return role.Msg.Role, nil
	} else if !errors.As(err, &cerr) || cerr.Code() != connect.CodeNotFound {
		return nil, err
	}

	role, err = cli.GetRole(root.Context(), connect.NewRequest(&idmv1.GetRoleRequest{
		Search: &idmv1.GetRoleRequest_Name{
			Name: roleNameOrId,
		},
	}))

	if err == nil {
		return role.Msg.Role, nil
	}

	if errors.As(err, &cerr) && cerr.Code() == connect.CodeNotFound {
		return nil, nil
	}

	return nil, err
}
