package models

type SearchQuery struct {
	Query string `json:"query" form:"query"`
	Limit int    `json:"limit" form:"limit"`
}

type SearchResult struct {
    Scores []DocumentScore `json:"scores"`
    Response    string `json:"response"`
}
