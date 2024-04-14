package controllers

import (
	"io"
	"log"
	"net/http"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentsController interface {
	GetDocument(c *gin.Context)
	GetDocuments(c *gin.Context)
	CreateDocument(c *gin.Context)
}

type documentsController struct {
	repo repositories.DocumentsRepository
}

func NewDocumentsController(repo repositories.DocumentsRepository) DocumentsController {
	return documentsController{repo}
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

    document, err := this.repo.GetDocument(c, slug, id)
    if err != nil {
		c.Status(http.StatusNotFound)
		return
    }

	c.JSON(http.StatusOK, document)
}

func (this documentsController) GetDocuments(c *gin.Context) {
	slug, err := uuid.Parse(c.Param("slug"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	documents, err := this.repo.GetDocuments(c, slug)
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

        document, err = this.repo.CreateDocument(c, document)
        if err != nil {
            log.Printf("Could not create document entry for file %s: %s", file.Filename, err.Error())
            continue
        }

        documents = append(documents, document)
	}

	c.JSON(http.StatusOK, documents)
}
