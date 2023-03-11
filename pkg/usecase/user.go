package usecase

import (
	"context"
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/repository/interfaces"
	services "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo   interfaces.UserRepository
	jwtUseCase services.JWTUseCase
}

func NewUserUseCase(userRepo interfaces.UserRepository, jwtUseCase services.JWTUseCase) services.UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		jwtUseCase: jwtUseCase,
	}
}

func (c *userUseCase) UserSignup(ctx context.Context, input domain.User) (domain.User, error) {
	//Hashing user password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password : %w", err)
	}
	input.Password = string(hash)
	userData, err := c.userRepo.CreateUser(ctx, input)

	return userData, err
}

func (c *userUseCase) UserLogin(ctx context.Context, user domain.User) (domain.User, string, string, error) {
	userDB, err := c.userRepo.FindUser(ctx, user.Email)
	if err != nil {
		return domain.User{}, "", "", err
	}
	if userDB.Email == "" {
		return domain.User{}, "", "", fmt.Errorf("user not found")
	}

	fmt.Println("User password - usecase : ", user.Password)
	fmt.Println("User DB password - usecase : ", userDB.Password)
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, "", "", fmt.Errorf("incorrect password: %w", err)
	}

	// generate access token
	accessToken, err := c.jwtUseCase.GenerateAccessToken(int(userDB.ID), userDB.Email, "user")
	if err != nil {
		return domain.User{}, "", "", fmt.Errorf("failed to create access token : %w", err)
	}

	// generate refresh token
	refreshToken, err := c.jwtUseCase.GenerateRefreshToken(int(userDB.ID), userDB.Email, "user")
	if err != nil {
		return domain.User{}, "", "", fmt.Errorf("failed to create refresh token : %w", err)
	}
	return userDB, accessToken, refreshToken, nil
}

func (c *userUseCase) ViewAllProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := c.userRepo.ViewAllProducts(ctx)
	return products, err
}
