package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PostDTO struct {
	Name        string `json:"name"`
	Description string `json:"description" binding:"required"`
}

type Post struct {
	bun.BaseModel `bun:"table:posts,alias:p"`

	Slug        uuid.UUID `bun:"slug,pk,type:uuid,default:uuid_generate_v4()" json:"slug"`
	Name        string    `bun:"name,type:varchar(128),notnull" json:"name"`
	Description string    `bun:"description,type:varchar(512),nullzero,notnull,default:''" json:"description"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	AuthorID    uuid.UUID `bun:"author_id,type:uuid,notnull" json:"authorId"`

    Author      *User     `bun:"rel:belongs-to,join:author_id=id" json:"author"`
    Documents []*Document `bun:"rel:has-many,join:slug=post_slug"`
}

func NewPost(authorId uuid.UUID, p PostDTO) Post {
	return Post{Name: p.Name, Description: p.Description, AuthorID: authorId}
}
