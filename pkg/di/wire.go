//go:build wireinject
// +build wireinject

package di

import (
	http "github.com/amalmadhu06/go-jwt-bearer-token/pkg/api"
	handler "github.com/amalmadhu06/go-jwt-bearer-token/pkg/api/handler"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/api/middleware"
	config "github.com/amalmadhu06/go-jwt-bearer-token/pkg/config"
	db "github.com/amalmadhu06/go-jwt-bearer-token/pkg/db"
	repository "github.com/amalmadhu06/go-jwt-bearer-token/pkg/repository"
	usecase "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {

	wire.Build(
		//database connection
		db.ConnectDatabase,

		//handler
		handler.NewAdminHandler,
		handler.NewUserHandler,

		//middleware
		//middleware.NewUserMiddleware,
		middleware.NewAdminMiddleware,
		//usecase
		usecase.NewAdminUseCase,
		usecase.NewUserUseCase,
		usecase.NewJWTUserService,

		//repository
		repository.NewAdminRepository,
		repository.NewUserRepository,

		// server connection
		http.NewServerHTTP,
	)
	return &http.ServerHTTP{}, nil

}
