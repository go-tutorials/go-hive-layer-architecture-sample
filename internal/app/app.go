package app

import (
	"context"

	"github.com/beltran/gohive"
	"github.com/core-go/health"
	h "github.com/core-go/hive/health"
	"github.com/core-go/log"

	"go-service/internal/user"
)

type ApplicationContext struct {
	Health *health.Handler
	User   user.UserTransport
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	configuration := gohive.NewConnectConfiguration()
	configuration.Database = "masterdata"
	connection, errConn := gohive.Connect(cfg.Hive.Host, cfg.Hive.Port, cfg.Hive.Auth, configuration)
	if errConn != nil {
		return nil, errConn
	}

	logError := log.LogError

	userHandler, err := user.NewUserHandler(connection, logError)
	if err != nil {
		return nil, err
	}

	hiveChecker := h.NewHealthChecker(connection)
	healthHandler := health.NewHandler(hiveChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
