// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/api"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/api/handler"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/api/middleware"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/config"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/db"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/repository"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase"
)

// Injectors from wire.go:

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gormDB)
	jwtUseCase := usecase.NewJWTUserService()
	userUseCase := usecase.NewUserUseCase(userRepository, jwtUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	adminRepository := repository.NewAdminRepository(gormDB)
	adminUseCase := usecase.NewAdminUseCase(adminRepository, jwtUseCase)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	middlewareMiddleware := middleware.NewAdminMiddleware(jwtUseCase)
	serverHTTP := http.NewServerHTTP(userHandler, adminHandler, middlewareMiddleware)
	return serverHTTP, nil
}
