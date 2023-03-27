package handler

import (
	"errors"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/common/response"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	services "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authHeaderKey = "Authorization"
)

type UserHandler struct {
	userUseCase services.UserUseCase
	jwtUseCase  services.JWTUseCase
}

func NewUserHandler(userUseCase services.UserUseCase, jwtUseCase services.JWTUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		jwtUseCase:  jwtUseCase,
	}
}

// UserSignup
// @Summary User Signup
// @ID user-signup
// @Description User can sign up with email and password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_credentials body domain.User{} true "user credentials for creating new account"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /signup [post]
func (cr *UserHandler) UserSignup(c *gin.Context) {
	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		r := response.Response{
			StatusCode: 400,
			Message:    "failed to read request body",
			Error:      err.Error()}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	user, err := cr.userUseCase.UserSignup(c.Request.Context(), newUser)
	if err != nil {
		r := response.Response{
			StatusCode: 400,
			Message:    "signup failed",
			Error:      err.Error()}
		c.JSON(http.StatusBadRequest, r)
		return
	}
	r := response.Response{
		StatusCode: 201,
		Message:    "successfully signed up",
		Data:       user}
	c.JSON(http.StatusCreated, r)
}

// UserLogin
// @Summary User Login
// @ID user-login
// @Description User can log in with email and password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_credentials body domain.User true "user credentials for logging in"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /login [post]
func (cr *UserHandler) UserLogin(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		r := response.Response{StatusCode: 400, Message: "failed to read request body", Error: err.Error()}
		c.JSON(http.StatusBadRequest, r)
		return
	}
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

// UserLogout
// @Summary User Logout
// @ID user-logout
// @Description User can log out of the website
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 307
// @Failure 400 {object} response.Response
// @Router /logout [post]
func (cr *UserHandler) UserLogout(c *gin.Context) {

	//remove access and refresh token from header
	c.Header("access-token", "")
	c.Header("refresh-token", "")

	//remove access_token from cookies if it exists
	//c.SetCookie("access_token", "", -1, "", "", false, true)
	//c.SetCookie("refresh_token", "", -1, "", "", false, true)

	//r := response.Response{StatusCode: 200, Message: "successfully logged out"}
	//c.JSON(http.StatusOK, r)

	//redirect to login page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// GetAccessToken
// @Summary Get access token with refresh token
// @ID get-access-token
// @Description Access token can be generated with a valid refresh token
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /access [post]
func (cr *UserHandler) GetAccessToken(c *gin.Context) {
	// Get refresh token from request header
	authHeader := c.GetHeader(authHeaderKey)
	if authHeader == "" {
		err := errors.New("request does not contain refresh token token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Extract the access token from the header
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		err := errors.New("invalid access token format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshToken := bearerToken[1]
	ok, claims, err := cr.jwtUseCase.VerifyToken(refreshToken)
	if !ok || err != nil {
		err := errors.New("invalid or expired refresh token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, err := cr.jwtUseCase.GenerateAccessToken(claims.ID, claims.Email, "user")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	r := response.Response{StatusCode: 201, Message: "successfully generated access token"}
	c.Header("access-token", accessToken)
	c.Header("refresh-token", refreshToken)
	c.JSON(http.StatusOK, r)

}

// ViewAllProducts
// @Summary View Products
// @ID view-products
// @Description User can can see listed products after logging in
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [get]
func (cr *UserHandler) ViewAllProducts(c *gin.Context) {
	products, err := cr.userUseCase.ViewAllProducts(c.Request.Context())
	if err != nil {
		r := response.Response{StatusCode: 500, Message: "failed to fetch product details", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	r := response.Response{StatusCode: 200, Message: "successfully fetched all products ", Data: products}
	c.JSON(http.StatusOK, r)
}
