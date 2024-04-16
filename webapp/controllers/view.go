package controllers

import (
	"net/http"

	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ViewController interface {
	GetIndexPage(c *gin.Context)
	GetHomePage(c *gin.Context)
	GetUserPage(c *gin.Context)
	GetPostPage(c *gin.Context)
	GetCreatePostPage(c *gin.Context)
	SearchPost(c *gin.Context)
}

type viewController struct {
	postsRepo         repositories.PostsRepository
	usersRepo         repositories.UsersRepository
	documentsRepo     repositories.DocumentsRepository
	embeddingsService services.EmbeddingsService
}

func NewViewController(postsRepo repositories.PostsRepository, usersRepo repositories.UsersRepository, documentsRepo repositories.DocumentsRepository, embeddingsService services.EmbeddingsService) ViewController {
	return viewController{postsRepo, usersRepo, documentsRepo, embeddingsService}
}

func (this viewController) GetIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (this viewController) GetHomePage(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	user, err := this.usersRepo.GetUser(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	posts, err := this.postsRepo.GetPosts(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Posts": posts, "User": user,
	})
}

func (this viewController) GetUserPage(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	user, err := this.usersRepo.GetUser(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "user.html", gin.H{"User": user})
}

type PostPageGetParams struct {
    Slug string `uri:"slug" binding:"required,uuid"`
}

func (this viewController) GetPostPage(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	user, err := this.usersRepo.GetUser(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

    params := PostPageGetParams{}
    if err := c.ShouldBindUri(&params); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    }

	post, err := this.postsRepo.GetPost(c, uuid.MustParse(params.Slug))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post, "User": user,
	})
}

func (this viewController) GetCreatePostPage(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	user, err := this.usersRepo.GetUser(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "create.html", gin.H{"User": user})
}

type SearchPageParams struct {
    Slug string `uri:"slug" binding:"required,uuid"`
}

func (this viewController) SearchPost(c *gin.Context) {
    params := SearchPageParams{}
    if err := c.ShouldBindUri(&params); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    }

	query := models.SearchQuery{Limit: 3}
	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

    searchResult, err := this.embeddingsService.GetSearchResult(c, uuid.MustParse(params.Slug), query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filter := models.DocumentsFilter{IDs: searchResult.DocumentIDs}
	documents, err := this.documentsRepo.GetDocuments(c, uuid.MustParse(params.Slug), filter)

    c.HTML(http.StatusOK, "search", gin.H{"Documents": documents, "Response": searchResult.Response})
}
