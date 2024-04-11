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

    if (err != nil) {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts,
	})
}

func (p viewController) GetPostPage(c *gin.Context) {
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

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post,
	})
}
