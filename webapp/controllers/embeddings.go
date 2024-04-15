package controllers

import (
	"net/http"
	"strconv"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmbeddingsController interface {
	GetSimilarDocument(c *gin.Context)
}

type embeddingsController struct {
	documentsRepo     repositories.DocumentsRepository
	embeddingsService services.EmbeddingsService
}

func NewEmbeddingsController(documentsRepo repositories.DocumentsRepository, embeddingsService services.EmbeddingsService) EmbeddingsController {
	return embeddingsController{documentsRepo, embeddingsService}
}

func (this embeddingsController) GetSimilarDocument(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	query, ok := c.GetQuery("query")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	limit := 10
	limitStr, ok := c.GetQuery("limit")
	if ok {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit parameter must be an integer"})
			return
		}

		limit = l
	}

	embeddings, err := this.embeddingsService.GetSimilarities(c, slug, query, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	documents := []models.Document{}
	for _, e := range embeddings {
		document, err := this.documentsRepo.GetDocument(c, slug, e.DocumentID)
		if err != nil {
			continue
		}

		documents = append(documents, document)
	}

	c.JSON(http.StatusOK, documents)
}
