package models

type SearchQuery struct {
	Query string `json:"query" form:"query"`
	Limit int    `json:"limit" form:"limit"`
}

type SearchResult struct {
	Scores   []DocumentScore `json:"scores"`
	Response string          `json:"response"`
}

type DocumentSearchResult struct {
	Filename string  `json:"filename"`
	Score    float32 `json:"score"`
}

func NewDocumentSearchResult(filename string, score float32) DocumentSearchResult {
	return DocumentSearchResult{Filename: filename, Score: score}
}
