package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Post struct {
    bun.BaseModel `bun:"table:posts,alias:u"`

	Slug        uuid.UUID `bun:",pk"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
}
