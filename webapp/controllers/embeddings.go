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

func (this embeddingsController) GetSearchResult(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	query := models.SearchQuery{Limit: 3}
	if err := c.Bind(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	searchResult, err := this.embeddingsService.GetSearchResult(c, slug, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, searchResult)
}
