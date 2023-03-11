package interfaces

import (
	"context"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	FindUser(ctx context.Context, email string) (domain.User, error)
	ViewAllProducts(ctx context.Context) ([]domain.Product, error)
}
