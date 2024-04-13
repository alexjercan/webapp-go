package controllers

import (
	"fmt"
	"net/http"

	"webapp-go/webapp/config"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
}

type authController struct {
	cfg config.Config
    service services.AuthService
    repo repositories.UsersRepository
}

func NewAuthController(cfg config.Config, service services.AuthService, repo repositories.UsersRepository) AuthController {
	return authController{cfg, service, repo}
}

func (this authController) Login(c *gin.Context) {
	state, err := this.service.GenerateRandomState()
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

    c.Redirect(http.StatusTemporaryRedirect, this.service.AuthCodeURL(state))
}

func (this authController) Callback(c *gin.Context) {
    session := sessions.Default(c)
    if c.Query("state") != session.Get("state") {
        c.String(http.StatusBadRequest, "Invalid state parameter.")
        return
    }

    token, err := this.service.AccessToken(c.Query("code"))
    if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    session.Set("access_token", token.AccessToken)
    if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    // TODO: This should be done in some other place I guess

    accessToken := session.Get("access_token").(string)
    dto, err := this.service.UserInfo(accessToken)
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    user, err := this.repo.CreateUser(c, models.NewUser(dto))
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    fmt.Println(user)

    //

    c.Redirect(http.StatusTemporaryRedirect, "/user")
}
