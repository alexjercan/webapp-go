package controllers

import (
	"log/slog"
	"net/http"

	"webapp-go/webapp/config"
	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CallbackQuery struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
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
    slog.Debug("[Login] Start")

	// Generate a random state
	state, err := this.authService.GenerateRandomState()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Login] State generated", "state", state)

	// Save the state in the session
	session := sessions.Default(c)
    session.AddFlash(state)
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Login] State saved in session")

	c.Redirect(http.StatusTemporaryRedirect, this.authService.AuthCodeURL(state))
}

func (this authController) Callback(c *gin.Context) {
    slog.Debug("[Callback] Start")

	// Parse the query parameters
	query := CallbackQuery{}
	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

    slog.Debug("[Callback] Query parsed", "query", query)

	session := sessions.Default(c)
    slog.Debug("[Callback] Session retrieved")

	// Check the state parameter
	flashes := session.Flashes()
    if len(flashes) == 0 {
        slog.Debug("[Callback] No state in session")
        c.String(http.StatusBadRequest, "Invalid state parameter.")
        return
    }

    state := flashes[0].(string)
	if state != query.State {
        slog.Debug("[Callback] State mismatch", "state", state, "queryState", query.State)
		c.String(http.StatusBadRequest, "Invalid state parameter.")
		return
	}

    slog.Debug("[Callback] State checked")

	// Get the access token
	token, err := this.authService.AccessToken(query.Code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Callback] Token received", "token", token)

	// Get the user info
	info, err := this.authService.UserInfo(token.AccessToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Callback] User info received", "info", info)

	// Create or update the user
	user, err := this.usersService.CreateOrUpdateUser(c, info)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Callback] User created or updated", "user", user)

	// Save the user ID in the session (logged in)
	session.Set(middlewares.USER_ID_KEY, user.ID.String())
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Callback] User ID saved in session")

	c.Redirect(http.StatusTemporaryRedirect, "/home")
}

func (this authController) Logout(c *gin.Context) {
    slog.Debug("[Logout] Start")

	session := sessions.Default(c)

	// Clear the session
	session.Clear()
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[Logout] Session cleared")

	c.Status(http.StatusNoContent)
}

func (this authController) BearerToken(c *gin.Context) {
    slog.Debug("[BearerToken] Start")

	// Get the user ID from the session
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

    slog.Debug("[BearerToken] User ID from session", "userId", userId)

	// Get the user to check that it exists
	user, err := this.usersService.GetUser(c, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    slog.Debug("[BearerToken] User found in database", "user", user)

	token, err := this.bearerService.GenerateToken(user.ID)

    slog.Debug("[BearerToken] Token generated", "token", token)

	c.String(http.StatusOK, token)
}

func (this authController) GetUser(c *gin.Context) {
    slog.Debug("[GetUser] Start")

	// API endpoint needs to get the user ID from the context
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

    slog.Debug("[GetUser] User ID from context", "userId", userId)

	// Get the user
	user, err := this.usersService.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

    slog.Debug("[GetUser] User found in database", "user", user)

	c.JSON(http.StatusOK, user)
}
