package usecase

import (
	"context"
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	interfaces "github.com/amalmadhu06/go-jwt-bearer-token/pkg/repository/interfaces"
	services "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
)

type adminUseCase struct {
	adminRepo  interfaces.AdminRepository
	jwtUseCase services.JWTUseCase
}

func NewAdminUseCase(adminRepo interfaces.AdminRepository, jwtUseCase services.JWTUseCase) services.AdminUseCase {
	return &adminUseCase{
		adminRepo:  adminRepo,
		jwtUseCase: jwtUseCase,
	}
}

func (c *adminUseCase) AdminLogin(ctx context.Context, admin domain.Admin) (domain.Admin, string, string, error) {
	adminDB, err := c.adminRepo.FindAdmin(ctx, admin.Email)
	if err != nil {
		return domain.Admin{}, "", "", err
	}
	if adminDB.Email == "" {
		return domain.Admin{}, "", "", err
	}
	if admin.Password != adminDB.Password {
		return domain.Admin{}, "", "", fmt.Errorf("incorrect password")

	}
	accessToken, err := c.jwtUseCase.GenerateAccessToken(int(adminDB.ID), adminDB.Email, "admin")
	if err != nil {
		return domain.Admin{}, "", "", err
	}
	refreshToken, err := c.jwtUseCase.GenerateRefreshToken(int(adminDB.ID), adminDB.Email, "admin")
	if err != nil {
		return domain.Admin{}, "", "", err
	}
	return adminDB, accessToken, refreshToken, nil
}

func (c *adminUseCase) AddProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error) {
	addedProduct, err := c.adminRepo.AddProduct(ctx, newProduct)
	return addedProduct, err
}

func (c *adminUseCase) ViewAllProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := c.adminRepo.ViewAllProducts(ctx)
	return products, err
}
