package interfaces

import (
	"context"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
)

type UserUseCase interface {
	UserSignup(ctx context.Context, input domain.User) (domain.User, error)
	UserLogin(ctx context.Context, user domain.User) (domain.User, string, string, error)
	ViewAllProducts(ctx context.Context) ([]domain.Product, error)
}
