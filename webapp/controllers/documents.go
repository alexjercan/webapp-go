package controllers

import (
	"io"
	"log"
	"net/http"
	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentsController interface {
	GetDocument(c *gin.Context)
	GetDocuments(c *gin.Context)
	CreateDocument(c *gin.Context)
	UpdateDocument(c *gin.Context)
	DeleteDocument(c *gin.Context)
}

type documentsController struct {
	documentsRepo repositories.DocumentsRepository
	postsRepo     repositories.PostsRepository
	documentChan  chan<- models.DocumentChanItem
}

func NewDocumentsController(documentsRepo repositories.DocumentsRepository, postsRepo repositories.PostsRepository, documentChan chan<- models.DocumentChanItem) DocumentsController {
	return documentsController{documentsRepo, postsRepo, documentChan}
}

func (this documentsController) GetDocument(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	document, err := this.documentsRepo.GetDocument(c, slug, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (this documentsController) GetDocuments(c *gin.Context) {
	var filter models.DocumentsFilter
	if err := c.Bind(&filter); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	documents, err := this.documentsRepo.GetDocuments(c, slug, filter)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (this documentsController) CreateDocument(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := this.postsRepo.GetPost(c, slug)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if post.AuthorID != userId.(uuid.UUID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["file"]

	documents := make([]models.Document, 0)
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			log.Printf("Could not open file %s", file.Filename)
			continue
		}

		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			log.Printf("Could not read file %s", file.Filename)
			continue
		}

		document := models.NewDocument(
			models.DocumentDTO{
				Filename:    file.Filename,
				ContentType: file.Header.Get("Content-Type"),
				Content:     content,
				PostSlug:    slug,
			},
		)

		document, err = this.documentsRepo.CreateDocument(c, document)
		if err != nil {
			log.Printf("Could not create document entry for file %s: %s", file.Filename, err.Error())
			continue
		}

		this.documentChan <- models.NewDocumentChanItem(models.CREATE, document.PostSlug, document.ID)

		documents = append(documents, document)
	}

	c.JSON(http.StatusOK, documents)
}

func (this documentsController) UpdateDocument(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := this.postsRepo.GetPost(c, slug)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if post.AuthorID != userId.(uuid.UUID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto models.DocumentDTO
	if err := c.Bind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	document, err := this.documentsRepo.UpdateDocument(c, slug, id, models.NewDocument(dto))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	this.documentChan <- models.NewDocumentChanItem(models.UPDATE, document.PostSlug, document.ID)

	c.JSON(http.StatusOK, document)
}

func (this documentsController) DeleteDocument(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId, exists := c.Get(middlewares.USER_ID_KEY)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := this.postsRepo.GetPost(c, slug)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if post.AuthorID != userId.(uuid.UUID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = this.documentsRepo.DeleteDocument(c, slug, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	this.documentChan <- models.NewDocumentChanItem(models.DELETE, slug, id)

	c.Status(http.StatusNoContent)
}
