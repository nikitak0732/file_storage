package usecase

import (
	"mime/multipart"
	"test_astral/internal/controller/DTO/request"
	"test_astral/internal/controller/DTO/response"
	"uuid"
)

type Cache interface {
	GetDocs(limit int) *response.DocsListResponse
	GetDocsById(id string) *response.Document
	InsertData(doc *response.Document)
	DeleteData(id string)
}
type User interface {
	CreateUser(login string, password string) (string, error)
	Auth(login string, password string) (string, error)
	GetAdminToken() string
	DeleteToken(token string) (string, error)
	GetSessionByToken(token string) (string, error)
}
type File interface {
	UploadDocument(
		ownerID uuid.UUID,
		meta request.UploadMeta,
		jsonData interface{},
		file multipart.File,
	) (*response.UploadData, *response.Document, error)
	ListDocuments(limit int) (*response.DocsListResponse, error)
	GetDocById(id string) (*response.Document, error)
	DeleteDocById(id string) error
}
