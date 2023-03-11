package interfaces

import (
	"context"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
)

type AdminRepository interface {
	FindAdmin(ctx context.Context, email string) (domain.Admin, error)
	AddProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error)
	ViewAllProducts(ctx context.Context) ([]domain.Product, error)
}
