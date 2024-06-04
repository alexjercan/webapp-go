package models

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DocumentDTO struct {
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	Content     []byte    `json:"content"`
	PostSlug    uuid.UUID `json:"postSlug"`
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
	switch this.ContentType {
	case "text/plain":
		fallthrough
	case "text/x-rst":
		fallthrough
	case "text/markdown":
		return string(this.Content)
	default:
		log.Printf("Content type not supported: %s\n", this.ContentType)
		return ""
	}
}

func (this Document) FormatPrompt() string {
	return fmt.Sprintf("Document Title: %s\nContent: %s\n", this.Filename, this.ParseContent())
}
