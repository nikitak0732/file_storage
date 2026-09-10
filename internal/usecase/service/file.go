package service

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"test_astral/internal/controller/DTO/request"
	"test_astral/internal/controller/DTO/response"
	"test_astral/internal/repo"
	"uuid"
)

type FileService struct {
	db repo.FileRepo
}

func NewFileService(db repo.FileRepo) *FileService {
	return &FileService{db: db}
}

func (s *FileService) UploadDocument(
	ownerID uuid.UUID,
	meta request.UploadMeta,
	jsonData interface{},
	file multipart.File,
) (*response.UploadData, *response.Document, error) {

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	var jsonBytes []byte
	if jsonData != nil {
		jsonBytes, err = json.Marshal(jsonData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	doc := response.Document{
		OwnerID:  ownerID,
		Name:     meta.Name,
		Mime:     meta.Mime,
		IsFile:   meta.File,
		IsPublic: meta.Public,
		JSONData: jsonBytes,
		FileData: fileBytes,
	}
	err = s.db.CreateDocument(doc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save document: %w", err)
	}

	fmt.Println("Создали документ")
	// 5. Возвращаем результат
	return &response.UploadData{
		JSON: jsonData,
		File: meta.Name,
	}, &doc, nil
}

func (s *FileService) ListDocuments(limit int) (*response.DocsListResponse, error) {
	docs, err := s.db.ListDocuments(limit)
	if err != nil {
		return nil, fmt.Errorf("получение списка документов: %w", err)
	}

	return &response.DocsListResponse{
		Data: response.DocsListData{
			Docs: docs,
		},
	}, nil

}
func (s *FileService) GetDocById(id string) (*response.Document, error) {
	return s.db.GetDocById(id)

}

func (s *FileService) DeleteDocById(id string) error {
	return s.db.DeleteDocById(id)
}
