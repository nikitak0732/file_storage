package response

import (
	"encoding/json"
	"time"
	"uuid"
)

type APIResponse struct {
	Error    *ErrorDetail    `json:"error,omitempty"`
	Response *ResponseDetail `json:"response,omitempty"`
	Data     *DataDetail     `json:"data,omitempty"`
}

type ErrorDetail struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type ResponseDetail struct {
	Login string `json:"login,omitempty"`
	Token string `json:"token,omitempty"`
}

type DataDetail struct {
	// UserID   int    `json:"user_id,omitempty"`
	// Username string `json:"username,omitempty"`
	// Email    string `json:"email,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

type UploadResponse struct {
	Data UploadData `json:"data"`
}

type UploadData struct {
	JSON interface{} `json:"json,omitempty"`
	File string      `json:"file"`
}
type DynamicResponse struct {
	Response map[string]bool `json:"response"`
}

func NewDynamicResponse(key string, value bool) *DynamicResponse {
	return &DynamicResponse{
		Response: map[string]bool{
			key: value,
		},
	}
}

type DocsListData struct {
	Docs []Document `json:"docs"`
}

type DocsListResponse struct {
	Data DocsListData `json:"data"`
}

type JSONResponse struct {
	Data json.RawMessage `json:"data"`
}

type Document struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	Name      string
	Mime      string
	IsFile    bool
	IsPublic  bool
	JSONData  []byte // JSONB
	FileData  []byte // BYTEA
	CreatedAt time.Time
	UpdatedAt time.Time
}
