package controllers

import (
	"net/http"

	"webapp-go/webapp/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ViewController interface {
	GetIndexPage(c *gin.Context)
	GetPostPage(c *gin.Context)
}

type viewController struct {
	db *bun.DB
}

func NewViewController(db *bun.DB) ViewController {
	return viewController{db}
}

func (p viewController) GetIndexPage(c *gin.Context) {
    var posts []models.Post
    err := p.db.NewSelect().Model(&posts).Scan(c)

    if (err != nil) {}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts,
	})
}

func (p viewController) GetPostPage(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

    post := new(models.Post)
    err = p.db.NewSelect().Model(post).Where("slug = ?", slug).Scan(c)

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post,
	})
}
