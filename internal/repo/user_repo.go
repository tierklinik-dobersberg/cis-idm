package repo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/rqlite/gorqlite"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo/stmts"
)

type UserRepo struct {
	Conn *gorqlite.Connection
}

func New(endpoint string) (*UserRepo, error) {
	conn, err := gorqlite.Open(endpoint)
	if err != nil {
		return nil, err
	}

	return &UserRepo{Conn: conn}, nil
}

func (repo *UserRepo) Migrate(ctx context.Context) error {
	if err := stmts.CreateUserTable.Write(ctx, repo.Conn, nil); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	if err := stmts.CreateAddressTable.Write(ctx, repo.Conn, nil); err != nil {
		return fmt.Errorf("failed to create user_addresses table: %w", err)
	}

	if err := stmts.CreatePhoneNumberTable.Write(ctx, repo.Conn, nil); err != nil {
		return fmt.Errorf("failed to create user_phone_numbers table: %w", err)
	}

	if err := stmts.CreateEMailTable.Write(ctx, repo.Conn, nil); err != nil {
		return fmt.Errorf("failed to create user_emails table: %w", err)
	}

	if err := stmts.CreateGroupTable.Write(ctx, repo.Conn, nil); err != nil {
		return fmt.Errorf("failed to create groups table: %w", err)
	}

	if err := stmts.CreateGroupMembershipTable.Write(ctx, repo.Conn, nil); err != nil {
		return fmt.Errorf("failed to create group_memberships table: %w", err)
	}
	return nil
}

func (repo *UserRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	return QueryOne(ctx, stmts.GetUserByID, repo.Conn, map[string]any{"id": id})
}

func (repo *UserRepo) GetUserByName(ctx context.Context, name string) (models.User, error) {
	return QueryOne(ctx, stmts.GetUserByName, repo.Conn, map[string]any{"username": name})
}

func (repo *UserRepo) GetGroupByName(ctx context.Context, name string) (models.Group, error) {
	return QueryOne(ctx, stmts.GetGroupByName, repo.Conn, map[string]any{
		"name": name,
	})
}

func (repo *UserRepo) GetGroupByID(ctx context.Context, id string) (models.Group, error) {
	return QueryOne(ctx, stmts.GetGroupByID, repo.Conn, map[string]any{
		"id": id,
	})
}

func (repo *UserRepo) CreateGroup(ctx context.Context, group models.Group) (models.Group, error) {
	if group.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return models.Group{}, err
		}

		group.ID = id.String()
	}

	if err := stmts.CreateGroup.Write(ctx, repo.Conn, group); err != nil {
		return models.Group{}, err
	}

	return group, nil
}

func (repo *UserRepo) AddGroupMembership(ctx context.Context, userID string, groupID string) error {
	return stmts.AddGroupMembership.Write(ctx, repo.Conn, models.GroupMembership{
		UserID:  userID,
		GroupID: groupID,
	})
}

func (repo *UserRepo) GetUsersInGroup(ctx context.Context, groupID string) ([]models.User, error) {
	return Query(ctx, stmts.GetUsersInGroup, repo.Conn, models.GroupMembership{
		GroupID: groupID,
	})
}

func (repo *UserRepo) GetUserGroupMemberships(ctx context.Context, userID string) ([]models.Group, error) {
	return Query(ctx, stmts.GetUserGroupMemberships, repo.Conn, models.GroupMembership{
		UserID: userID,
	})
}

func (repo *UserRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	if user.ID == "" {
		userID, err := uuid.NewV4()
		if err != nil {
			return user, err
		}

		user.ID = userID.String()
	}

	if err := stmts.CreateUser.Write(ctx, repo.Conn, user); err != nil {
		return user, err
	}

	return user, nil
}

func Query[T any](ctx context.Context, stmt stmts.Statement[T], conn *gorqlite.Connection, args any) ([]T, error) {
	pStmt, err := stmt.Prepare(args)
	if err != nil {
		return nil, err
	}

	queryResult, err := conn.QueryOneParameterizedContext(ctx, pStmt)
	if err != nil {
		if queryResult.Err != nil {
			return nil, queryResult.Err
		}

		return nil, err
	}

	typeOf := reflect.TypeOf(stmt.Result)
	results := make([]T, 0, queryResult.NumRows())

	for queryResult.Next() {
		m, err := queryResult.Map()
		if err != nil {
			return results, err
		}

		obj := reflect.New(typeOf).Interface().(*T)
		if err := mapstructure.Decode(m, obj); err != nil {
			return results, err
		}

		results = append(results, *obj)
	}

	return results, nil
}

func QueryOne[T any](ctx context.Context, stmt stmts.Statement[T], conn *gorqlite.Connection, args any) (T, error) {
	results, err := Query(ctx, stmt, conn, args)
	if err != nil {
		return stmt.Result, err
	}

	if len(results) == 0 {
		return stmt.Result, stmts.ErrNoResults
	}

	if len(results) > 1 {
		return stmt.Result, fmt.Errorf("query returned more than one result")
	}

	return results[0], nil
}
