package repository

import (
	"context"
	"fmt"

	"github.com/beltran/gohive"

	"go-service/internal/user/model"
)

type UserAdapter struct {
	Connection *gohive.Connection
}

func NewUserRepository(connection *gohive.Connection) *UserAdapter {
	return &UserAdapter{Connection: connection}
}

func (m *UserAdapter) All(ctx context.Context) ([]model.User, error) {
	cursor := m.Connection.Cursor()
	query := "select id, username, email, phone, status, createdDate from users"
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	var res []model.User
	var user model.User
	for cursor.HasMore(ctx) {
		cursor.FetchOne(ctx, &user.Id, &user.Username, &user.Email, &user.Phone, &user.Status, &user.CreatedDate)
		if cursor.Err != nil {
			return nil, cursor.Err
		}

		res = append(res, user)
	}
	return res, nil
}

func (m *UserAdapter) Load(ctx context.Context, id string) (*model.User, error) {
	cursor := m.Connection.Cursor()
	var user model.User
	query := fmt.Sprintf("select id, username, email, phone, status , createdDate from users where id = '%v' order by id asc limit 1", id)

	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	for cursor.HasMore(ctx) {
		cursor.FetchOne(ctx, &user.Id, &user.Username, &user.Email, &user.Phone, &user.Status, &user.CreatedDate)
		if cursor.Err != nil {
			return nil, cursor.Err
		}
		return &user, nil
	}
	return nil, nil
}

func (m *UserAdapter) Create(ctx context.Context, user *model.User) (int64, error) {
	cursor := m.Connection.Cursor()
	query := fmt.Sprintf("insert into users values ('%v', '%v', '%v', '%v', '%v', '%v')", user.Id, user.Username, user.Email, user.Phone, user.Status, user.CreatedDate)
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}

func (m *UserAdapter) Update(ctx context.Context, user *model.User) (int64, error) {
	cursor := m.Connection.Cursor()
	query := fmt.Sprintf("update users set username = '%v', email = '%v', phone = '%v' where id = '%v'", user.Username, user.Email, user.Phone, user.Id)
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}

func (m *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
	cursor := m.Connection.Cursor()
	query := fmt.Sprintf("delete from users where id = '%v'", id)
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	return 1, nil
}
