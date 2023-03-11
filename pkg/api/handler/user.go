package handler

import (
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/common/response"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	services "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(userUseCase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (cr *UserHandler) UserSignup(c *gin.Context) {
	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		r := response.Response{StatusCode: 422, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	}
	user, err := cr.userUseCase.UserSignup(c.Request.Context(), newUser)
	if err != nil {
		r := response.Response{StatusCode: 400, Message: "signup failed", Error: err.Error()}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	r := response.Response{StatusCode: 201, Message: "successfully signed up", Data: user}
	c.JSON(http.StatusCreated, r)
}

func (cr *UserHandler) UserLogin(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		r := response.Response{StatusCode: 422, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusUnprocessableEntity, r)
		return
	}
	fmt.Println("user credentials : handler", user)
	loggedInUser, accessToken, refreshToken, err := cr.userUseCase.UserLogin(c.Request.Context(), user)
	if err != nil {
		r := response.Response{StatusCode: 400, Message: "failed to login as user", Error: err.Error()}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	r := response.Response{StatusCode: 200, Message: "successfully logged in as user", Data: loggedInUser}
	c.Header("access-token", accessToken)
	c.Header("refresh-token", refreshToken)
	c.JSON(http.StatusOK, r)
}

func (cr *UserHandler) UserLogout(c *gin.Context) {
	c.Header("access-token", "")
	c.Header("refresh-token", "")
	r := response.Response{StatusCode: 200, Message: "successfully logged out"}
	c.JSON(http.StatusOK, r)
}

func (cr *UserHandler) ViewAllProducts(c *gin.Context) {
	products, err := cr.userUseCase.ViewAllProducts(c.Request.Context())
	if err != nil {
		r := response.Response{StatusCode: 500, Message: "failed to fetch product details", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	r := response.Response{StatusCode: 302, Message: "successfully fetched all products ", Data: products}
	c.JSON(http.StatusFound, r)
}
