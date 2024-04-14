package services

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
)

type AuthService interface {
	AuthCodeURL(state string) string
	AccessToken(code string) (models.TokenResponse, error)
	UserInfo(accessToken string) (models.GitHubUser, error)
	GenerateRandomState() (string, error)
}

type githubAuthService struct {
	cfg config.Config
}

func NewAuthService(cfg config.Config) AuthService {
	return githubAuthService{cfg}
}

func (this githubAuthService) GenerateRandomState() (state string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}

	state = base64.StdEncoding.EncodeToString(b)

	return
}

func (this githubAuthService) AuthCodeURL(state string) string {
	values := url.Values{}
	values.Add("client_id", this.cfg.OAuth.ClientId)
	values.Add("scope", "user")
	values.Add("redirect_uri", this.cfg.OAuth.RedirectUri)
	values.Add("state", state)
	query := values.Encode()

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", query)
}

func (this githubAuthService) accessTokenURL(code string) string {
	values := url.Values{}
	values.Add("client_id", this.cfg.OAuth.ClientId)
	values.Add("client_secret", this.cfg.OAuth.ClientSecret)
	values.Add("code", code)
	query := values.Encode()

	return fmt.Sprintf("https://github.com/login/oauth/access_token?%s", query)
}

func (this githubAuthService) AccessToken(code string) (token models.TokenResponse, err error) {
	url := this.accessTokenURL(code)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &token)

	return
}

func (this githubAuthService) UserInfo(accessToken string) (user models.GitHubUser, err error) {
	url := "https://api.github.com/user"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	user = models.GitHubUser{}
	err = json.Unmarshal(body, &user)

	return
}
