package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DocumentCommand int

const (
	CREATE DocumentCommand = iota
	UPDATE
	DELETE
)

type DocumentChanItem struct {
	Command DocumentCommand
    PostSlug uuid.UUID
	ID      uuid.UUID
}

func NewDocumentChanItem(c DocumentCommand, slug uuid.UUID, id uuid.UUID) DocumentChanItem {
    return DocumentChanItem{Command: c, PostSlug: slug, ID: id}
}

type DocumentEmbedding struct {
	bun.BaseModel `bun:"table:document_embeddings,alias:de"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
    DocumentID  uuid.UUID `bun:"document_id,type:uuid,notnull,unique" json:"documentId"`
    Embeddings  []float32 `bun:"embeddings,type:vector(4096),notnull" json:"embeddings"`

    Document    *Document `bun:"rel:has-one,join:document_id=id" json:"document"`
}

func NewDocumentEmbedding(id uuid.UUID, embeddings []float32) DocumentEmbedding {
    return DocumentEmbedding{DocumentID: id, Embeddings: embeddings}
}
