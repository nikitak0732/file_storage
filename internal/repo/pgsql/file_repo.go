package repo

import (
	"database/sql"
	"fmt"
	"test_astral/internal/controller/DTO/response"
	"test_astral/pkg/postgres"
	"uuid"
)

type FileRepo struct {
	pg *postgres.Postgres
}

func NewFileRepo(pg *postgres.Postgres) *FileRepo {
	return &FileRepo{pg}
}
func (f *FileRepo) CreateDocument(doc response.Document) error {
	query := `
        INSERT INTO documents (
            owner_id,
            name,
            mime,
            is_file,
            is_public,
            json_data,
            file_data
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
       
    `
	if len(doc.JSONData) == 0 {
		doc.JSONData = []byte(`{}`)
	}
	if len(doc.FileData) == 0 {
		doc.FileData = []byte(`{}`)
	}
	err := f.pg.Postgres.QueryRow(query,
		uuid.UUID(doc.OwnerID),
		doc.Name,
		doc.Mime,
		doc.IsFile,
		doc.IsPublic,
		doc.JSONData,
		doc.FileData,
	)

	if err.Err() != nil {
		return fmt.Errorf("Ошибка при добавлении документа: %s", err.Err().Error())
	}

	return nil
}
func (f *FileRepo) ListDocuments(limit int) ([]response.Document, error) {
	query := `
        SELECT
            id,
            owner_id,
            name,
            mime,
            is_file,
            is_public,
            json_data,
            file_data,
            created_at,
            updated_at
        FROM documents
        ORDER BY name ASC, created_at DESC
      
    `
	var (
		rows *sql.Rows
		err  error
	)
	if limit != 0 {
		query += "LIMIT $1"
		rows, err = f.pg.Postgres.Query(query, limit)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения списка документов: %w", err)
		}
	} else {
		rows, err = f.pg.Postgres.Query(query)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения списка документов: %w", err)
		}
	}

	defer rows.Close()

	var docs []response.Document
	for rows.Next() {
		var doc response.Document
		if err := rows.Scan(
			&doc.ID,
			&doc.OwnerID,
			&doc.Name,
			&doc.Mime,
			&doc.IsFile,
			&doc.IsPublic,
			&doc.JSONData,
			&doc.FileData,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ошибка сканирования документа: %w", err)
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %w", err)
	}

	if docs == nil {
		docs = []response.Document{}
	}
	return docs, nil
}

func (f *FileRepo) GetDocById(id string) (*response.Document, error) {
	query := `
        SELECT
            id,
            owner_id,
            name,
            mime,
            is_file,
            is_public,
            json_data,
            file_data,
            created_at,
            updated_at
        FROM documents
        WHERE id = $1::uuid
      
       
    `
	var doc response.Document
	err := f.pg.Postgres.QueryRow(query, string(id)).Scan(
		&doc.ID,
		&doc.OwnerID,
		&doc.Name,
		&doc.Mime,
		&doc.IsFile,
		&doc.IsPublic,
		&doc.JSONData,
		&doc.FileData,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения документов: %w", err)
	}

	return &doc, nil
}

func (f *FileRepo) DeleteDocById(id string) error {
	query := `
	delete from documents where id = $1
    `
	_, err := f.pg.Postgres.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления документов: %w", err)
	}

	return nil
}
