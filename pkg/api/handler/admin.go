package handler

import (
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/common/response"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	services "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(adminUseCase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: adminUseCase,
	}
}

func (cr *AdminHandler) AdminLogin(c *gin.Context) {
	var admin domain.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		r := response.Response{StatusCode: 422, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	}
	loggedInAdmin, accessToken, refreshToekn, err := cr.adminUseCase.AdminLogin(c.Request.Context(), admin)
	if err != nil {
		r := response.Response{StatusCode: 400, Message: "failed to login as admin", Error: err.Error()}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	r := response.Response{StatusCode: 200, Message: "successfully logged in as admin", Data: loggedInAdmin}
	c.Header("access-token", accessToken)
	c.Header("refresh-token", refreshToekn)
	c.JSON(http.StatusOK, r)

}

func (cr *AdminHandler) AdminLogout(c *gin.Context) {
	c.Header("access-token", "")
	c.Header("refresh-token", "")
	r := response.Response{StatusCode: 200, Message: "successfully logged out"}
	c.JSON(http.StatusOK, r)
}

//For simplicity, adding 'Add Product' method handler here. It should ideally be in seperate product file

func (cr *AdminHandler) AddProduct(c *gin.Context) {
	var newProduct domain.Product

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		r := response.Response{StatusCode: 422, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	}
	product, err := cr.adminUseCase.AddProduct(c.Request.Context(), newProduct)
	if err != nil {
		r := response.Response{StatusCode: 500, Message: "failed to add new product", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	r := response.Response{StatusCode: 201, Message: "successfully created new product", Data: product}
	c.JSON(http.StatusCreated, r)
}

func (cr *AdminHandler) ViewAllProducts(c *gin.Context) {
	products, err := cr.adminUseCase.ViewAllProducts(c.Request.Context())
	if err != nil {
		r := response.Response{StatusCode: 500, Message: "failed to fetch product details", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	r := response.Response{StatusCode: 302, Message: "successfully created new product", Data: products}
	c.JSON(http.StatusFound, r)
}
