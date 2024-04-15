package services

import (
	"context"
	"log"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms/ollama"
)

type EmbeddingsService interface {
	GetSimilarities(c context.Context, slug uuid.UUID, query string, limit int) ([]models.DocumentEmbedding, error)
	Worker(c context.Context)
}

type embeddingsService struct {
	documentsRepo  repositories.DocumentsRepository
	embeddingsRepo repositories.EmbeddingsRepository
	llm            *ollama.LLM
	documentChan   <-chan models.DocumentChanItem
}

func NewEmbeddingsService(documentsRepo repositories.DocumentsRepository, embeddingsRepo repositories.EmbeddingsRepository, llm *ollama.LLM, documentChan <-chan models.DocumentChanItem) EmbeddingsService {
	return embeddingsService{documentsRepo, embeddingsRepo, llm, documentChan}
}

func (this embeddingsService) createEmbeddings(c context.Context, slug uuid.UUID, id uuid.UUID) {
	document, err := this.documentsRepo.GetDocument(c, slug, id)
	if err != nil {
		log.Printf("Error getting the document with id %s: %s\n", id, err.Error())
		return
	}

	content := document.ParseContent()

	embeddings, err := this.llm.CreateEmbedding(c, []string{content})
	if err != nil {
		log.Printf("Error generating embeddings for document with id %s: %s\n", id, err.Error())
		return
	}

	_, err = this.embeddingsRepo.CreateEmbedding(c, models.NewDocumentEmbedding(id, embeddings[0]))
	if err != nil {
		log.Printf("Error saving the embeddings for document with id %s: %s\n", id, err.Error())
		return
	}
}

func (this embeddingsService) updateEmbeddings(c context.Context, slug uuid.UUID, id uuid.UUID) {
	document, err := this.documentsRepo.GetDocument(c, slug, id)
	if err != nil {
		log.Printf("Error getting the document with id %s: %s\n", id, err.Error())
		return
	}

	content := document.ParseContent()

	embeddings, err := this.llm.CreateEmbedding(c, []string{content})
	if err != nil {
		log.Printf("Error generating embeddings for document with id %s: %s\n", id, err.Error())
		return
	}

	_, err = this.embeddingsRepo.UpdateEmbedding(c, id, models.NewDocumentEmbedding(id, embeddings[0]))
	if err != nil {
		log.Printf("Error saving the embeddings for document with id %s: %s\n", id, err.Error())
		return
	}
}

func (this embeddingsService) deleteEmbeddings(c context.Context, id uuid.UUID) {
	_, err := this.embeddingsRepo.DeleteEmbedding(c, id)
	if err != nil {
		log.Printf("Error deleting the embeddings for document with id %s: %s\n", id, err.Error())
		return
	}
}

func (this embeddingsService) Worker(c context.Context) {
	for d := range this.documentChan {
		switch d.Command {
		case models.CREATE:
			this.createEmbeddings(c, d.PostSlug, d.ID)
		case models.UPDATE:
			this.updateEmbeddings(c, d.PostSlug, d.ID)
		case models.DELETE:
			this.deleteEmbeddings(c, d.ID)
		default:
			log.Printf("Unkown command: %d\n", d.Command)
		}
	}
}

func (this embeddingsService) GetSimilarities(c context.Context, slug uuid.UUID, query string, limit int) (embeddings []models.DocumentEmbedding, err error) {
	embeddings = []models.DocumentEmbedding{}

	es, err := this.llm.CreateEmbedding(c, []string{query})
	if err != nil {
		return
	}

	embeddings, err = this.embeddingsRepo.GetSimilarEmbeddings(c, slug, es[0], limit)

	return
}
