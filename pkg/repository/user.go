package repository

import (
	"context"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userRepo{DB}
}

func (c *userRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	var userDB domain.User
	createQuery := `INSERT INTO users(email, password) VALUES($1,$2) RETURNING *;`
	err := c.DB.Raw(createQuery, user.Email, user.Password).Scan(&userDB).Error
	return userDB, err
}

func (c *userRepo) FindUser(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	fetchQuery := `SELECT * FROM users WHERE email = $1`
	err := c.DB.Raw(fetchQuery, email).Scan(&user).Error
	return user, err
}

func (c *userRepo) ViewAllProducts(ctx context.Context) ([]domain.Product, error) {
	var allProducts []domain.Product
	findAllQuery := `SELECT * FROM products;`
	rows, err := c.DB.Raw(findAllQuery).Rows()
	if err != nil {
		return allProducts, err
	}

	defer rows.Close()
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return allProducts, err
		}
		allProducts = append(allProducts, product)
	}

	return allProducts, err
}
