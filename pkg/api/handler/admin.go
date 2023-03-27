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

// AdminLogin
// @Summary Admin Login
// @ID admin-login
// @Description Admin can log in with email and password
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param admin_credentials body domain.Admin true "admin credentials for logging in"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin-panel/login [post]
func (cr *AdminHandler) AdminLogin(c *gin.Context) {
	var admin domain.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		r := response.Response{StatusCode: 400, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusBadRequest, r)
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

// AdminLogout
// @Summary Admin Logout
// @ID admin-logout
// @Description Admin can log out of the website
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Product json
// @Success 200
// @Router /admin-panel/logout [post]
func (cr *AdminHandler) AdminLogout(c *gin.Context) {
	c.Header("access-token", "")
	c.Header("refresh-token", "")
	r := response.Response{StatusCode: 200, Message: "successfully logged out"}
	c.JSON(http.StatusOK, r)
	c.Redirect(http.StatusTemporaryRedirect, "/adminpanel")
}

//For simplicity, adding 'Add Product' method handler here. It should ideally be in separate product file

// AddProduct
// @Summary Add product to inventory
// @ID add-product
// @Description Admin can add products to inventory
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Product json
// @Param product_details body domain.Product true "product details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin-panel/product [post]
func (cr *AdminHandler) AddProduct(c *gin.Context) {
	var newProduct domain.Product

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		r := response.Response{StatusCode: 400, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusBadRequest, r)
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

// ViewAllProducts
// @Summary View Products
// @ID admin-view-products
// @Description Admin can can see listed products
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin-panel/product [get]
func (cr *AdminHandler) ViewAllProducts(c *gin.Context) {
	products, err := cr.adminUseCase.ViewAllProducts(c.Request.Context())
	if err != nil {
		r := response.Response{StatusCode: 500, Message: "failed to fetch product details", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	r := response.Response{StatusCode: 200, Message: "successfully created new product", Data: products}
	c.JSON(http.StatusOK, r)
}
