package adapter

import (
	"context"
	"fmt"
	"reflect"

	"github.com/beltran/gohive"
	h "github.com/core-go/hive"

	"go-service/internal/user/model"
)

type UserAdapter struct {
	Connection    *gohive.Connection
	Map           map[string]int
	Keys          []string
	JsonColumnMap map[string]string
}

func NewUserRepository(connection *gohive.Connection) *UserAdapter {
	userType := reflect.TypeOf(model.User{})
	fieldsIndex, _, jsonColumnMap, keys, _, _, _ := h.Init(userType)
	return &UserAdapter{Connection: connection, Map: fieldsIndex, Keys: keys, JsonColumnMap: jsonColumnMap}
}

func (m *UserAdapter) All(ctx context.Context) ([]model.User, error) {
	cursor := m.Connection.Cursor()
	query := "select id, username, email, phone, date_of_birth from users"
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	var res []model.User

	for cursor.HasMore(ctx) {
		var user model.User
		var dob string
		cursor.FetchOne(ctx, &user.Id, &user.Username, &user.Email, &user.Phone, &dob)
		if cursor.Err != nil {
			return nil, cursor.Err
		}
		user.DateOfBirth = h.GetTime(dob)

		res = append(res, user)
	}
	return res, nil
}

func (m *UserAdapter) Load(ctx context.Context, id string) (*model.User, error) {
	cursor := m.Connection.Cursor()

	query := fmt.Sprintf("select id, username, email, phone, date_of_birth from users where id = '%s' order by id asc limit 1", h.Escape(id))

	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	for cursor.HasMore(ctx) {
		var user model.User
		var dob string
		cursor.FetchOne(ctx, &user.Id, &user.Username, &user.Email, &user.Phone, &dob)
		if cursor.Err != nil {
			return nil, cursor.Err
		}
		user.DateOfBirth = h.GetTime(dob)
		return &user, nil
	}
	return nil, nil
}

func (m *UserAdapter) Create(ctx context.Context, user *model.User) (int64, error) {
	cursor := m.Connection.Cursor()
	query := fmt.Sprintf("insert into users(id, username, email, phone, date_of_birth) values ('%s', '%s', %s, '%s', %s)", h.Escape(user.Id), h.Escape(user.Username), h.GetString(user.Email), h.Escape(user.Phone), h.GetDateTime(user.DateOfBirth))
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}

func (m *UserAdapter) Update(ctx context.Context, user *model.User) (int64, error) {
	cursor := m.Connection.Cursor()
	query := fmt.Sprintf("update users set username = '%s', email = %s, phone = '%s', date_of_birth= %s where id = '%s'", h.Escape(user.Username), h.GetString(user.Email), h.Escape(user.Phone), h.GetDateTime(user.DateOfBirth), h.Escape(user.Id))
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}

func (m *UserAdapter) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	colMap := h.JSONToColumns(user, m.JsonColumnMap)
	query := h.BuildToPatch("users", colMap, m.Keys)
	cursor := m.Connection.Cursor()
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}

func (m *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
	cursor := m.Connection.Cursor()
	query := fmt.Sprintf("delete from users where id = '%s'", h.Escape(id))
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}
