package repositories

import (
	"context"
	"webapp-go/webapp/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EmbeddingsRepository interface {
	GetEmbedding(c context.Context, id uuid.UUID) (models.DocumentEmbedding, error)
	GetEmbeddingFor(c context.Context, documentID uuid.UUID) (models.DocumentEmbedding, error)
	GetSimilarEmbeddings(c context.Context, slug uuid.UUID, embedding []float32, limit int) ([]models.DocumentEmbedding, error)
	CreateEmbedding(c context.Context, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error)
	UpdateEmbedding(c context.Context, id uuid.UUID, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error)
	DeleteEmbedding(c context.Context, id uuid.UUID) (uuid.UUID, error)
}

type embeddingsRepository struct {
	db *bun.DB
}

func NewEmbeddingsRepository(db *bun.DB) EmbeddingsRepository {
	return embeddingsRepository{db}
}

func (this embeddingsRepository) GetEmbedding(c context.Context, id uuid.UUID) (embedding models.DocumentEmbedding, err error) {
	err = this.db.NewSelect().Model(&embedding).Where("id = ?", id).Scan(c)

	return
}

func (this embeddingsRepository) GetEmbeddingFor(c context.Context, documentID uuid.UUID) (embedding models.DocumentEmbedding, err error) {
	err = this.db.NewSelect().Model(&embedding).Where("document_id = ?", documentID).Scan(c)

	return
}

func (this embeddingsRepository) GetSimilarEmbeddings(c context.Context, slug uuid.UUID, embedding []float32, limit int) ([]models.DocumentEmbedding, error) {
	embeddings := []models.DocumentEmbedding{}

	err := this.db.NewSelect().Model(&embeddings).Relation("Document").Where("post_slug = ?", slug).OrderExpr("embeddings <-> ?", embedding).Limit(limit).Scan(c)

	return embeddings, err
}

func (this embeddingsRepository) CreateEmbedding(c context.Context, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error) {
	_, err := this.db.NewInsert().Model(&embedding).Exec(c)

	return embedding, err
}

func (this embeddingsRepository) UpdateEmbedding(c context.Context, id uuid.UUID, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error) {
	embedding.ID = id

	_, err := this.db.NewUpdate().Model(&embedding).OmitZero().WherePK().Exec(c)

	return embedding, err
}

func (this embeddingsRepository) DeleteEmbedding(c context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, err := this.db.NewDelete().Model(&models.DocumentEmbedding{}).Where("id = ?", id).Exec(c)

	return id, err
}
