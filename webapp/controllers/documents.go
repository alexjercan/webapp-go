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

type DocumentGetQuery struct {
	Slug uuid.UUID `form:"slug" binding:"required"`
	ID   uuid.UUID `form:"id" binding:"required"`
}

func (this documentsController) GetDocument(c *gin.Context) {
	query := DocumentGetQuery{}
	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	document, err := this.documentsRepo.GetDocument(c, query.Slug, query.ID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, document)
}

type DocumentsGetQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
}

func (this documentsController) GetDocuments(c *gin.Context) {
	query := DocumentsGetQuery{}
	if err := c.ShouldBindUri(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	documents, err := this.documentsRepo.GetDocuments(c, uuid.MustParse(query.Slug))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

type DocumentCreateQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
}

func (this documentsController) CreateDocument(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	query := DocumentCreateQuery{}
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
				PostSlug:    post.Slug,
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

type DocumentUpdateQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
	ID   string `uri:"id" binding:"required,uuid"`
}

func (this documentsController) UpdateDocument(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	query := DocumentUpdateQuery{}
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

	var dto models.DocumentDTO
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	document, err := this.documentsRepo.UpdateDocument(c, post.Slug, uuid.MustParse(query.ID), models.NewDocument(dto))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	this.documentChan <- models.NewDocumentChanItem(models.UPDATE, document.PostSlug, document.ID)

	c.JSON(http.StatusOK, document)
}

type DocumentDeleteQuery struct {
	Slug string `uri:"slug" binding:"required,uuid"`
	ID   string `uri:"id" binding:"required,uuid"`
}

func (this documentsController) DeleteDocument(c *gin.Context) {
	userId := c.MustGet(middlewares.USER_ID_KEY).(uuid.UUID)

	query := DocumentDeleteQuery{}
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

	_, err = this.documentsRepo.DeleteDocument(c, post.Slug, uuid.MustParse(query.ID))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	this.documentChan <- models.NewDocumentChanItem(models.DELETE, post.Slug, uuid.MustParse(query.ID))

	c.Status(http.StatusNoContent)
}
