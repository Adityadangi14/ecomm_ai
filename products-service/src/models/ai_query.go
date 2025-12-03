package models

type AiQueryParams struct {
	Query     string `json:"query"`
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`
	OrgID     string `json:"orgId"`
}

type MessageChanStruct struct {
	Chunk string
	Err   error
}
