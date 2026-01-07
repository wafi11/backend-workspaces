package types

type PaginationCursor struct {
	NextCursor *string `json:"nextCursor,omitempty"`
	HasMore    bool    `json:"hasMore"`
	Count      int     `json:"count"`
}
