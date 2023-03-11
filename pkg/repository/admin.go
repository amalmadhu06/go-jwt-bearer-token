package repository

import (
	"context"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type adminRepo struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepo{DB}
}

func (c *adminRepo) FindAdmin(ctx context.Context, email string) (domain.Admin, error) {
	var admin domain.Admin
	fetchQuery := `SELECT * FROM admins WHERE email = $1;`
	err := c.DB.Raw(fetchQuery, email).Scan(&admin).Error
	return admin, err
}

func (c *adminRepo) AddProduct(ctx context.Context, newProduct domain.Product) (domain.Product, error) {
	var addedProduct domain.Product
	addQuery := `INSERT INTO products(name, description, price) VALUES($1, $2, $3) RETURNING *;`
	err := c.DB.Raw(addQuery, newProduct.Name, newProduct.Description, newProduct.Price).Scan(&addedProduct).Error
	return addedProduct, err
}

func (c *adminRepo) ViewAllProducts(ctx context.Context) ([]domain.Product, error) {
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