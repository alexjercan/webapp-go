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
	postsRepo repositories.PostsRepository
    usersRepo repositories.UsersRepository
}

func NewPostsController(postsRepo repositories.PostsRepository, usersRepo repositories.UsersRepository) PostsController {
	return postsController{postsRepo, usersRepo}
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

	post, err := this.postsRepo.GetPost(c, uuid.MustParse(query.Slug))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (this postsController) GetPosts(c *gin.Context) {
	posts, err := this.postsRepo.GetPosts(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (this postsController) CreatePost(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

    user, err := this.usersRepo.GetUser(c, userId)
    if err != nil {
        c.AbortWithError(http.StatusNotFound, err)
        return
    }

    if user.IsAnonymous() {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

	dto := models.PostDTO{}
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post, err := this.postsRepo.CreatePost(c, models.NewPost(userId, dto))
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

	post, err := this.postsRepo.GetPost(c, uuid.MustParse(query.Slug))
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

	post, err = this.postsRepo.UpdatePost(c, post.Slug, models.NewPost(userId, dto))
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

	post, err := this.postsRepo.GetPost(c, uuid.MustParse(query.Slug))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if post.AuthorID != userId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = this.postsRepo.DeletePost(c, post.Slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
