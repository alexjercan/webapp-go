package controllers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"webapp-go/webapp/config"
	"webapp-go/webapp/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
}

type authController struct {
	cfg config.Config
}

func NewAuthController(cfg config.Config) AuthController {
	return authController{cfg}
}

func generateRandomState() (state string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}

	state = base64.StdEncoding.EncodeToString(b)

	return
}

func (p authController) authCodeURL(state string) string {
    values := url.Values{}
    values.Add("client_id", p.cfg.OAuth.ClientId)
    values.Add("scope", "user")
    values.Add("redirect_uri", p.cfg.OAuth.RedirectUri)
    values.Add("state", state)
    query := values.Encode()

    return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", query)
}

func (p authController) accessTokenURL(code string) string {
    values := url.Values{}
    values.Add("client_id", p.cfg.OAuth.ClientId)
    values.Add("client_secret", p.cfg.OAuth.ClientSecret)
    values.Add("code", code)
    query := values.Encode()

    return fmt.Sprintf("https://github.com/login/oauth/access_token?%s", query)
}

func (p authController) Login(c *gin.Context) {
	state, err := generateRandomState()
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

    c.Redirect(http.StatusTemporaryRedirect, p.authCodeURL(state))
}

func (p authController) Callback(c *gin.Context) {
    session := sessions.Default(c)
    if c.Query("state") != session.Get("state") {
        c.String(http.StatusBadRequest, "Invalid state parameter.")
        return
    }

    url := p.accessTokenURL(c.Query("code"))

    client := &http.Client{}
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }
    req.Header.Add("Accept", "application/json")

    res, err := client.Do(req)
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    token := models.TokenResponse{}
    err = json.Unmarshal(body, &token)
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.Redirect(http.StatusTemporaryRedirect, "/")
}
