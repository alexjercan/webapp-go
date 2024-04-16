package models

import "github.com/google/uuid"

type SearchQuery struct {
	Query string `json:"query" form:"query"`
	Limit int    `json:"limit" form:"limit"`
}

type SearchResult struct {
    DocumentIDs []uuid.UUID `json:"documentIDs"`
    Response    string `json:"response"`
}
