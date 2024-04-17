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

type PostGetQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
}

func (this postsController) GetPost(c *gin.Context) {
	query := PostGetQuery{}
	if err := c.ShouldBindUri(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := this.repo.GetPost(c, uuid.MustParse(query.Slug))
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
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	dto := models.PostDTO{}
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := this.repo.CreatePost(c, models.NewPost(userId, dto))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

type PostUpdateQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
}

func (this postsController) UpdatePost(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	query := PostUpdateQuery{}
	if err := c.ShouldBindUri(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := this.repo.GetPost(c, uuid.MustParse(query.Slug))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if post.AuthorID != userId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dto := models.PostDTO{}
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err = this.repo.UpdatePost(c, post.Slug, models.NewPost(userId, dto))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

type PostDeleteQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
}

func (this postsController) DeletePost(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	query := PostDeleteQuery{}
	if err := c.ShouldBindUri(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := this.repo.GetPost(c, uuid.MustParse(query.Slug))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if post.AuthorID != userId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = this.repo.DeletePost(c, post.Slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
