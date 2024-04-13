package controllers

import (
	"net/http"

	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"

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
    service services.AuthService
}

func NewViewController(repo repositories.PostsRepository, service services.AuthService) ViewController {
	return viewController{repo, service}
}

func (this viewController) GetIndexPage(c *gin.Context) {
    posts, err := this.repo.GetPosts(c)
    if (err != nil) {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts,
	})
}

func (this viewController) GetUserPage(c *gin.Context) {
    session := sessions.Default(c)
    accessToken := session.Get("access_token").(string)

    user, err := this.service.UserInfo(accessToken)
    if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.HTML(http.StatusOK, "user.html", gin.H{ "User": user })
}

func (this viewController) GetPostPage(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
        return
	}

    post, err := this.repo.GetPost(c, slug)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post,
	})
}
