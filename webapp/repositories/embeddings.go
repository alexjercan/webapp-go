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
	GetSimilarEmbeddings(c context.Context, slug uuid.UUID, embedding []float32, limit int) ([]models.DocumentScore, error)
	CreateEmbedding(c context.Context, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error)
	UpdateEmbeddingFor(c context.Context, documentID uuid.UUID, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error)
	DeleteEmbeddingFor(c context.Context, documentID uuid.UUID) (uuid.UUID, error)
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

func (this embeddingsRepository) GetSimilarEmbeddings(c context.Context, slug uuid.UUID, embedding []float32, limit int) ([]models.DocumentScore, error) {
	scores := []models.DocumentScore{}

	err := this.db.NewSelect().
		Table("document_embeddings").
		Column("document_embeddings.document_id").
		ColumnExpr("1 - (embeddings <=> ?) AS score", embedding).
		Join("JOIN documents as d").
		JoinOn("document_embeddings.document_id = d.id").
		Where("post_slug = ?", slug).
		Order("score").
		Limit(limit).
		Scan(c, &scores)

	return scores, err
}

func (this embeddingsRepository) CreateEmbedding(c context.Context, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error) {
	_, err := this.db.NewInsert().Model(&embedding).Exec(c)

	return embedding, err
}

func (this embeddingsRepository) UpdateEmbeddingFor(c context.Context, documentID uuid.UUID, embedding models.DocumentEmbedding) (models.DocumentEmbedding, error) {
	_, err := this.db.NewUpdate().Model(&embedding).OmitZero().Where("document_id = ?", documentID).Exec(c)

	return embedding, err
}

func (this embeddingsRepository) DeleteEmbeddingFor(c context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, err := this.db.NewDelete().Model(&models.DocumentEmbedding{}).Where("document_id = ?", id).Exec(c)

	return id, err
}
