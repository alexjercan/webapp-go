package repositories

import (
	"context"
	"webapp-go/webapp/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DocumentsRepository interface {
    GetDocument(c context.Context, slug uuid.UUID, id uuid.UUID) (models.Document, error)
    GetDocuments(c context.Context, slug uuid.UUID) ([]models.Document, error)
	CreateDocument(c context.Context, document models.Document) (models.Document, error)
}

type documentsRepository struct {
	db *bun.DB
}

func NewDocumentsRepository(db *bun.DB) DocumentsRepository {
	return documentsRepository{db}
}

func (this documentsRepository) GetDocument(c context.Context, slug uuid.UUID, id uuid.UUID) (document models.Document, err error) {
	err = this.db.NewSelect().Model(&document).Where("post_slug = ?", slug).Where("id = ?", id).Scan(c);

    return
}

func (this documentsRepository) GetDocuments(c context.Context, slug uuid.UUID) (documents []models.Document, err error) {
	err = this.db.NewSelect().Model(&documents).Where("post_slug = ?", slug).Scan(c)

    return
}

func (this documentsRepository) CreateDocument(c context.Context, document models.Document) (models.Document, error) {
	_, err := this.db.NewInsert().Model(&document).Exec(c)

    return document, err
}
