package service

import (
	"log"
	"sync"
	"test_astral/internal/controller/DTO/response"
	repo "test_astral/internal/repo/pgsql"
)

type Cache struct {
	Cache map[string]response.Document
	Mu    sync.RWMutex // Допустим у нас будет сильно больше чтений чем записей
}

func NewCache(db *repo.FileRepo) *Cache {
	docs, err := db.ListDocuments(0)
	if err != nil {
		log.Println("Не удалось загрузить кэш. Данные будут начитываться с БД")
	}
	mp := make(map[string]response.Document, len(docs))
	for _, doc := range docs {
		mp[doc.ID.String()] = doc
	}
	log.Println("Кэш загрузился")
	return &Cache{Mu: sync.RWMutex{}, Cache: mp}
}

func (c Cache) GetDocs(limit int) *response.DocsListResponse {

	c.Mu.RLock()
	defer c.Mu.RUnlock()
	if limit <= 0 {
		limit = len(c.Cache)
	}
	docs := make([]response.Document, 0, min(limit, len(c.Cache)))
	for _, doc := range c.Cache {
		if len(docs) >= limit {
			break
		}
		docs = append(docs, doc)
	}
	log.Println("Достали документы из кэша")
	return &response.DocsListResponse{
		Data: response.DocsListData{
			Docs: docs,
		},
	}
}

func (c Cache) GetDocsById(id string) *response.Document {

	c.Mu.RLock()
	defer c.Mu.RUnlock()
	log.Println("Достали документ из кэша")
	docs := c.Cache[id]
	return &docs
}

func (c *Cache) InsertData(doc *response.Document) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	log.Println("Докинули кэш")
	c.Cache[doc.ID.String()] = *doc
}

func (c *Cache) DeleteData(id string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.Cache, id)
}
