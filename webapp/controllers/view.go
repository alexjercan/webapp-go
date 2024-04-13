package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ViewController interface {
	GetIndexPage(c *gin.Context)
    GetUserPage(c *gin.Context)
	GetPostPage(c *gin.Context)
}

type viewController struct {
	repo repositories.PostsRepository
}

func NewViewController(repo repositories.PostsRepository) ViewController {
	return viewController{repo}
}

func (p viewController) GetIndexPage(c *gin.Context) {
    posts, err := p.repo.GetPosts(c)
    if (err != nil) {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts,
	})
}

func (p viewController) GetUserPage(c *gin.Context) {
    session := sessions.Default(c)
    accessToken := session.Get("access_token").(string)

    url := "https://api.github.com/user"

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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

    user := models.GitHubUser{}
    err = json.Unmarshal(body, &user)
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.HTML(http.StatusOK, "user.html", gin.H{ "User": user })
}

func (p viewController) GetPostPage(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
        return
	}

    post, err := p.repo.GetPost(c, slug)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post,
	})
}
