package controllers

import (
	"net/http"

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

func (p postsController) GetPost(c *gin.Context) {
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

	c.JSON(http.StatusOK, post)
}

func (p postsController) GetPosts(c *gin.Context) {
	posts, err := p.repo.GetPosts(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (p postsController) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := p.repo.CreatePost(c, post)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (p postsController) UpdatePost(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post models.Post
	if err := c.Bind(&post); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = p.repo.UpdatePost(c, slug, post)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (p postsController) DeletePost(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = p.repo.DeletePost(c, slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
