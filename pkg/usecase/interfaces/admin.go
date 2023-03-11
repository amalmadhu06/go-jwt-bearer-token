package interfaces

import (
	"context"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
)

type AdminUseCase interface {
	AdminLogin(ctx context.Context, admin domain.Admin) (domain.Admin, string, string, error)
	AddProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error)
	ViewAllProducts(ctx context.Context) ([]domain.Product, error)
}
