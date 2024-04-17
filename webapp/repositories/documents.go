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
	UpdateDocument(c context.Context, slug uuid.UUID, id uuid.UUID, document models.Document) (models.Document, error)
	DeleteDocument(c context.Context, slug uuid.UUID, id uuid.UUID) (uuid.UUID, error)
}

type documentsRepository struct {
	db *bun.DB
}

func NewDocumentsRepository(db *bun.DB) DocumentsRepository {
	return documentsRepository{db}
}

func (this documentsRepository) GetDocument(c context.Context, slug uuid.UUID, id uuid.UUID) (document models.Document, err error) {
	err = this.db.NewSelect().Model(&document).Where("post_slug = ?", slug).Where("id = ?", id).Scan(c)

	return
}

func (this documentsRepository) GetDocuments(c context.Context, slug uuid.UUID) (documents []models.Document, err error) {
	documents = []models.Document{}

	q := this.db.NewSelect().Model(&documents).Where("post_slug = ?", slug)

	err = q.Scan(c)

	return
}

func (this documentsRepository) CreateDocument(c context.Context, document models.Document) (models.Document, error) {
	_, err := this.db.NewInsert().Model(&document).Exec(c)

	return document, err
}

func (this documentsRepository) UpdateDocument(c context.Context, slug uuid.UUID, id uuid.UUID, document models.Document) (models.Document, error) {
	document.PostSlug = slug
	document.ID = id

	_, err := this.db.NewUpdate().Model(&document).OmitZero().WherePK().Exec(c)

	return document, err
}

func (this documentsRepository) DeleteDocument(c context.Context, slug uuid.UUID, id uuid.UUID) (uuid.UUID, error) {
	_, err := this.db.NewDelete().Model(&models.Document{}).Where("post_slug = ?", slug).Where("id = ?", id).Exec(c)

	return id, err
}
