package controllers

import (
	"net/http"

	"webapp-go/webapp/config"
	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
	BearerToken(c *gin.Context)
	GetUser(c *gin.Context)
}

type authController struct {
	cfg           config.Config
	authService   services.AuthService
	usersService  services.UsersService
	bearerService services.BearerService
}

func NewAuthController(cfg config.Config, authService services.AuthService, usersService services.UsersService, bearerService services.BearerService) AuthController {
	return authController{cfg, authService, usersService, bearerService}
}

func (this authController) Login(c *gin.Context) {
	state, err := this.authService.GenerateRandomState()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(c)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, this.authService.AuthCodeURL(state))
}

func (this authController) Callback(c *gin.Context) {
	session := sessions.Default(c)
	if c.Query("state") != session.Get("state") {
		c.String(http.StatusBadRequest, "Invalid state parameter.")
		return
	}

	token, err := this.authService.AccessToken(c.Query("code"))
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set("access_token", token.AccessToken)
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dto, err := this.authService.UserInfo(token.AccessToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := this.usersService.CreateOrUpdateUser(c, dto)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set(middlewares.USER_ID_KEY, user.ID.String())
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (this authController) BearerToken(c *gin.Context) {
	session := sessions.Default(c)
	accessToken := session.Get("access_token").(string)
	dto, err := this.authService.UserInfo(accessToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := this.usersService.GetUserByLogin(c, dto.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	t, err := this.bearerService.GenerateToken(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": t,
	})
}

func (this authController) GetUser(c *gin.Context) {
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := this.usersService.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
