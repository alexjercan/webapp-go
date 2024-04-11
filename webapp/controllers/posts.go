package controllers

import (
	"net/http"

	"webapp-go/webapp/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PostsController interface {
	GetPost(c *gin.Context)
	GetPosts(c *gin.Context)
	CreatePost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
}

type postsController struct {
	db *bun.DB
}

func NewPostsController(db *bun.DB) PostsController {
	return postsController{db}
}

func (p postsController) GetPost(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post := new(models.Post)
	if err := p.db.NewSelect().Model(post).Where("slug = ?", slug).Scan(c); err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (p postsController) GetPosts(c *gin.Context) {
	var posts []models.Post
	err := p.db.NewSelect().Model(&posts).Scan(c)

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

	_, err := p.db.NewInsert().Model(&post).Exec(c)
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

	post.Slug = slug

	_, err = p.db.NewUpdate().Model(&post).OmitZero().WherePK().Exec(c)
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

	_, err = p.db.NewDelete().Model((*models.Post)(nil)).Where("slug = ?", slug).Exec(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
