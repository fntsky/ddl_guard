package schema

import "time"

type CreateDraftReq struct {
	Raw   []byte    `json:"raw"`
	Type  string    `json:"data_type"`
	Draft DraftInfo `json:"draft"`
}

type DraftInfo struct {
	UUID        string    `json:"uuid"`
	Title       string    `json:"title"`
	Deadline    time.Time `json:"deadline"`
	EealyRemind bool      `json:"early_remind"`
	Description string    `json:"description"`
}
