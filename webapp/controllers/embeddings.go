package controllers

import (
	"net/http"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmbeddingsController interface {
	GetSearchResult(c *gin.Context)
}

type embeddingsController struct {
	documentsRepo     repositories.DocumentsRepository
	embeddingsService services.EmbeddingsService
}

func NewEmbeddingsController(documentsRepo repositories.DocumentsRepository, embeddingsService services.EmbeddingsService) EmbeddingsController {
	return embeddingsController{documentsRepo, embeddingsService}
}

type SearchGetParams struct {
	Slug string `uri:"slug" binding:"required,uuid"`
}

func (this embeddingsController) GetSearchResult(c *gin.Context) {
	params := SearchGetParams{}
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, searchResult)
}
