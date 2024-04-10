package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"webapp-go/webapp/models"
)

type PostsController interface {
	GetPost(c *gin.Context)
}

type postsController struct {
	db *bun.DB
}

func NewPostsController(db *bun.DB) PostsController {
	return postsController{db}
}

func (p postsController) GetPost(c *gin.Context) {
	slug := c.Param("slug")

	post := new(models.Post)
	if err := p.db.NewSelect().Model(post).Where("slug = ?", slug).Scan(c); err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, post)
}
