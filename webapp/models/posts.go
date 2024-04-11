package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:posts,alias:p"`

    Slug        uuid.UUID `bun:"slug,pk,type:uuid,default:uuid_generate_v4()" json:"slug"`
	Name        string    `bun:"name,type:varchar(128),notnull" json:"name" binding:"required"`
	Description string    `bun:"description,type:varchar(512),nullzero,notnull,default:''" json:"description"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
}
