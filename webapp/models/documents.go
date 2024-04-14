package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DocumentDTO struct {
	Filename    string
	ContentType string
	Content     []byte
	PostSlug    uuid.UUID
}

type Document struct {
	bun.BaseModel `bun:"table:documents,alias:d"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
    Filename    string    `bun:"filename,type:varchar(128),notnull,unique:post_group" json:"filename"`
	ContentType string    `bun:"content_type,type:varchar(128),notnull,default:'text/plain'" json:"contentType"`
	Content     []byte    `bun:"content,type:bytea,notnull,default:''" json:"content"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
    PostSlug    uuid.UUID `bun:"post_slug,type:uuid,notnull,unique:post_group" json:"postSlug"`
}

func NewDocument(d DocumentDTO) Document {
	return Document{Filename: d.Filename, ContentType: d.ContentType, Content: d.Content, PostSlug: d.PostSlug}
}

func (this Document) ParseContent() string {
    return string(this.Content)
}
