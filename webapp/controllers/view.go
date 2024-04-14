package controllers

import (
	"net/http"

	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ViewController interface {
	GetIndexPage(c *gin.Context)
	GetUserPage(c *gin.Context)
	GetPostPage(c *gin.Context)
	GetCreatePostPage(c *gin.Context)
}

type viewController struct {
	postsRepo repositories.PostsRepository
	usersRepo repositories.UsersRepository
}

func NewViewController(postsRepo repositories.PostsRepository, usersRepo repositories.UsersRepository) ViewController {
	return viewController{postsRepo, usersRepo}
}

func (this viewController) GetIndexPage(c *gin.Context) {
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := this.usersRepo.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	posts, err := this.postsRepo.GetPosts(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts, "User": user,
	})
}

func (this viewController) GetUserPage(c *gin.Context) {
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := this.usersRepo.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "user.html", gin.H{"User": user})
}

func (this viewController) GetPostPage(c *gin.Context) {
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := this.usersRepo.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := this.postsRepo.GetPost(c, slug)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post, "User": user,
	})
}

func (this viewController) GetCreatePostPage(c *gin.Context) {
	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := this.usersRepo.GetUser(c, userId.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "create.html", gin.H{"User": user})
}
