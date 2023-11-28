package cmds

import (
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
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

func resolveUser(root *cli.Root, userNameOrId string) (*idmv1.User, error) {
	cli := root.Users()

	user, err := cli.GetUser(root.Context(), connect.NewRequest(&idmv1.GetUserRequest{
		Search: &idmv1.GetUserRequest_Id{
			Id: userNameOrId,
		},
	}))

	var cerr *connect.Error
	if err == nil {
		return user.Msg.Profile.User, nil
	}

	if !errors.As(err, &cerr) || cerr.Code() != connect.CodeNotFound {
		return nil, err
	}

	user, err = cli.GetUser(root.Context(), connect.NewRequest(&idmv1.GetUserRequest{
		Search: &idmv1.GetUserRequest_Name{
			Name: userNameOrId,
		},
	}))

	if err == nil {
		return user.Msg.Profile.User, nil
	}

	if errors.As(err, &cerr) && cerr.Code() == connect.CodeNotFound {
		return nil, nil
	}

	return nil, err
}

func mustResolveUserToId(root *cli.Root, userNameOrId string) string {
	user, err := resolveUser(root, userNameOrId)
	if err != nil {
		logrus.Fatal(fmt.Errorf("failed to resolve user %q: %s", userNameOrId, err.Error()))
	}

	if user == nil {
		logrus.Fatal(fmt.Errorf("user not found %q", userNameOrId))
	}

	return user.Id
}

func mustResolveUserIds(root *cli.Root, userNamesOrIds []string) []string {
	result := make([]string, len(userNamesOrIds))

	for idx, u := range userNamesOrIds {
		id := mustResolveUserToId(root, u)
		result[idx] = id
	}

	return result
}
