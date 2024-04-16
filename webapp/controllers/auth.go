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

const STATE_KEY = "state"

type CallbackQuery struct {
    State string `form:"state" binding:"required"`
    Code string `form:"code" binding:"required"`
}

type AuthController interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
	Logout(c *gin.Context)
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
    // Generate a random state
	state, err := this.authService.GenerateRandomState()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    // Save the state in the session
	session := sessions.Default(c)
	session.Set(STATE_KEY, state)
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, this.authService.AuthCodeURL(state))
}

func (this authController) Callback(c *gin.Context) {
    // Parse the query parameters
    query := CallbackQuery{}
    if err := c.ShouldBind(&query); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    }

    // Check the state parameter
	session := sessions.Default(c)
    state := session.Get(STATE_KEY)
	if state != query.State {
		c.String(http.StatusBadRequest, "Invalid state parameter.")
		return
	}

    // Clear the state parameter
    session.Delete(STATE_KEY)
    if err := session.Save(); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    // Get the access token
	token, err := this.authService.AccessToken(query.Code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    // Get the user info
	info, err := this.authService.UserInfo(token.AccessToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    // Create or update the user
	user, err := this.usersService.CreateOrUpdateUser(c, info)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    // Save the user ID in the session (logged in)
	session.Set(middlewares.USER_ID_KEY, user.ID.String())
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/home")
}

func (this authController) Logout(c *gin.Context) {
	session := sessions.Default(c)

    // Clear the session
	session.Clear()
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (this authController) BearerToken(c *gin.Context) {
	session := sessions.Default(c)

    // Get the user ID from the session
	userId := uuid.MustParse(session.Get(middlewares.USER_ID_KEY).(string))

    // Get the user to check that it exists
	user, err := this.usersService.GetUser(c, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	token, err := this.bearerService.GenerateToken(user.ID)

	c.String(http.StatusOK, token)
}

func (this authController) GetUser(c *gin.Context) {
    // API endpoint needs to get the user ID from the context
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

    // Get the user
	user, err := this.usersService.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
