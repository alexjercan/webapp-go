package controllers

import (
	"net/http"

	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostsController interface {
	GetPost(c *gin.Context)
	GetPosts(c *gin.Context)
	CreatePost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
}

type postsController struct {
	repo repositories.PostsRepository
}

func NewPostsController(repo repositories.PostsRepository) PostsController {
	return postsController{repo}
}

func (this postsController) GetPost(c *gin.Context) {
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

	c.JSON(http.StatusOK, post)
}

func (this postsController) GetPosts(c *gin.Context) {
	posts, err := this.repo.GetPosts(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (this postsController) CreatePost(c *gin.Context) {
	var dto models.PostDTO
	if err := c.Bind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

    userId, exists := c.Get(middlewares.USER_ID_KEY)
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

	post, err := this.repo.CreatePost(c, models.NewPost(userId.(uuid.UUID), dto))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (this postsController) UpdatePost(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var dto models.PostDTO
	if err := c.Bind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

    userId, exists := c.Get(middlewares.USER_ID_KEY)
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    post, err := this.repo.GetPost(c, slug)
    if err != nil {
        c.Status(http.StatusNotFound)
        return
    }

    if post.AuthorID != userId.(uuid.UUID) {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    post, err = this.repo.UpdatePost(c, slug, models.NewPost(userId.(uuid.UUID), dto))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (this postsController) DeletePost(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

    userId, exists := c.Get(middlewares.USER_ID_KEY)
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    post, err := this.repo.GetPost(c, slug)
    if err != nil {
        c.Status(http.StatusNotFound)
        return
    }

    if post.AuthorID != userId.(uuid.UUID) {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

	_, err = this.repo.DeletePost(c, slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
