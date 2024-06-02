package user

import (
	"context"
	"net/http"
	"reflect"

	"github.com/beltran/gohive"
	v "github.com/core-go/core/v10"
	"github.com/core-go/hive"
	"github.com/core-go/search/hive/query"

	"go-service/internal/user/handler"
	"go-service/internal/user/model"
	"go-service/internal/user/repository"
	"go-service/internal/user/service"
)

type UserTransport interface {
	All(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(connection *gohive.Connection, logError func(context.Context, string, ...map[string]interface{})) (UserTransport, error) {
	validator, err := v.NewValidator()
	if err != nil {
		return nil, err
	}

	userType := reflect.TypeOf(model.User{})
	userQuery := query.NewBuilder("users", userType)
	userSearchBuilder, err := hive.NewSearchBuilder(connection, userType, userQuery.BuildQuery)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(connection)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userSearchBuilder.Search, userService, validator.Validate, logError)
	return userHandler, nil
}
