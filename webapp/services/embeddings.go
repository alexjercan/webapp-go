package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms/ollama"
)

type EmbeddingsService interface {
    GetSearchResult(c context.Context, slug uuid.UUID, query models.SearchQuery) (models.SearchResult, error)
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

func (this embeddingsService) updateEmbeddings(c context.Context, slug uuid.UUID, documentID uuid.UUID) {
	document, err := this.documentsRepo.GetDocument(c, slug, documentID)
	if err != nil {
		log.Printf("Error getting the document with id %s: %s\n", documentID, err.Error())
		return
	}

	content := document.ParseContent()

	embeddings, err := this.llm.CreateEmbedding(c, []string{content})
	if err != nil {
		log.Printf("Error generating embeddings for document with id %s: %s\n", documentID, err.Error())
		return
	}

	_, err = this.embeddingsRepo.UpdateEmbeddingFor(c, documentID, models.NewDocumentEmbedding(documentID, embeddings[0]))
	if err != nil {
		log.Printf("Error saving the embeddings for document with id %s: %s\n", documentID, err.Error())
		return
	}
}

func (this embeddingsService) deleteEmbeddings(c context.Context, documentID uuid.UUID) {
	_, err := this.embeddingsRepo.DeleteEmbeddingFor(c, documentID)
	if err != nil {
		log.Printf("Error deleting the embeddings for document with id %s: %s\n", documentID, err.Error())
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

func (this embeddingsService) buildPrompt(question string, documents []models.Document) string {
    prompt := "Based on the following context answer the given question as best as you can\n"

    contents := []string{}
    for _, d := range documents {
        contents = append(contents, d.ParseContent())
    }
    context := strings.Join(contents, "\n")

    return fmt.Sprintf("%s\n%s\nQuestion: %s\nAnswer: ", prompt, context, question)
}

func (this embeddingsService) GetSearchResult(c context.Context, slug uuid.UUID, query models.SearchQuery) (result models.SearchResult, err error) {
	es, err := this.llm.CreateEmbedding(c, []string{query.Query})
	if err != nil {
		return
	}

    scores, err := this.embeddingsRepo.GetSimilarEmbeddings(c, slug, es[0], query.Limit)
	if err != nil {
		return
	}

	documentIds := []uuid.UUID{}
	for _, s := range scores {
		documentIds = append(documentIds, s.DocumentID)
	}

    result.DocumentIDs = documentIds

	filter := models.DocumentsFilter{IDs: documentIds}
	documents, err := this.documentsRepo.GetDocuments(c, slug, filter)

    prompt := this.buildPrompt(query.Query, documents)

    response, err := this.llm.Call(c, prompt)
    if (err != nil) {
        return
    }

    result.Response = response

	return
}
