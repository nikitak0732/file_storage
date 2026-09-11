package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test_astral/internal/controller/DTO/request"
	"test_astral/internal/controller/DTO/response"
	"uuid"

	"github.com/go-chi/chi"
)

func (s *V1) RegisterRoutesFile(r chi.Router) {
	r.Post("/api/docs", s.UploadDocs)
	r.Get("/api/docs", s.GetDocs)
	r.Delete("/api/docs/{token}/{id}", s.DeleteDocs)
	r.Get("/api/docs/{token}/{id}", s.GetDocsById)
	r.Head("/api/docs", s.GetDocs)

	r.Head("/api/docs/{id}", s.GetDocsById)
	r.Delete("/api/docs/{id}", s.DeleteDocs)
}

// UploadDocs godoc
// @Summary      Загрузка документа
// @Description  Загружает документ
// @Tags         docs
// @Accept multipart/form-data
// @Produce      json
// @Param meta formData string true "JSON с параметрами запроса"
// @Param json formData string false "JSON с данными документа"
// @Param file formData file true "Файл документа"
// @Success      200  {object}  response.UploadData  "Файл загружен"
// @Failure      400  {object}  response.APIResponse  "Неверный запрос"
// @Failure      401  {object}  response.APIResponse  "Не авторизовались"
// @Failure      500  {object}  response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/docs [post]
func (s *V1) UploadDocs(w http.ResponseWriter, r *http.Request) {
	// Ограничимся на 32 метра
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Ошибка парсинга формы: %s", err.Error()))
		return
	}
	metaStr := r.FormValue("meta")
	if metaStr == "" {
		sendError(w, http.StatusBadRequest, "Поле 'meta' обязательно")
		return
	}
	var meta request.UploadMeta
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Ошибка парсинга meta: %s", err.Error()))
		return
	}

	if meta.Token == "" {
		sendError(w, http.StatusUnauthorized, "Токен обязателен")
		return
	}

	ownerID, err := s.U.GetSessionByToken(meta.Token)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Неверный токен")
		return
	}
	if err := validateMeta(meta); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Файл обязателен: %s", err.Error()))
		return
	}
	const maxFileSize = 10 << 20 // 10MB
	if fileHeader.Size > maxFileSize {
		sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Файл слишком большой (макс. %d MB)", maxFileSize/1024/1024))
		return
	}
	var jsonData any
	jsonStr := r.FormValue("json")
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			sendError(w, http.StatusBadRequest,
				fmt.Sprintf("Ошибка парсинга json: %s", err.Error()))
			return
		}
	}
	id, err := uuid.Parse(ownerID)
	if err != nil {
		sendError(w, http.StatusInternalServerError,
			fmt.Sprintf("Ошибка парсинга uuid: %s", err.Error()))
		return
	}
	result, doc, err := s.F.UploadDocument(
		id,
		meta,
		jsonData,
		file,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError,
			fmt.Sprintf("Ошибка вставки в БД: %s", err.Error()))
		return
	}
	s.C.InsertData(doc)
	resp := response.UploadResponse{
		Data: response.UploadData{
			JSON: result.JSON,
			File: result.File,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

// GetDocs godoc
// @Summary      Получение документов
// @Description  Получение документов
// @Tags         docs
// @Accept 		 json
// @Produce      json
// @Param        token  query     string  true   "Токен авторизации"
// @Param        limit  query     int     false  "Кол-во документов"
// @Success      200  {object}  response.DocsListResponse  "Получили документы"
// @Failure      400  {object}   response.APIResponse   "Неверный запрос"
// @Failure      401  {object}  response.APIResponse  "Не авторизовались"
// @Failure      500  {object}   response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/docs [get]
func (s *V1) GetDocs(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	token := q.Get("token")
	if token == "" {
		sendError(w, http.StatusUnauthorized, "Токен обязателен")
		return
	}

	if _, err := s.U.GetSessionByToken(token); err != nil {
		sendError(w, http.StatusUnauthorized, "Неверный токен")
		return
	}

	limit := 0
	if l := q.Get("limit"); l != "" {
		n, err := strconv.Atoi(l)
		if err != nil || n <= 0 {
			sendError(w, http.StatusBadRequest, "Неверный limit")
			return
		}
		if n > 100 {
			n = 100
		}
		limit = n
	}

	_, err := s.U.GetSessionByToken(token)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Неверный токен")
		return
	}
	var docs *response.DocsListResponse
	docs = s.C.GetDocs(limit)
	if len(docs.Data.Docs) == 0 {
		docs, err = s.F.ListDocuments(limit)
		if err != nil {
			sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка получения списков документов: %s", err.Error()))
			return
		}
	}

	resp := docs

	if docs == nil {
		resp = &response.DocsListResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetDocsById godoc
// @Summary      Получение документа по ID
// @Description  Получение документа по ID
// @Tags         docs
// @Accept 		json
// @Produce      json, image/jpeg, image/png, image/gif, image/webp
// @Param        token  query     string  true   "Токен авторизации"
// @Param        id  query     string     true  "ID документа"
// @Success      200  {file}  binary  "Изображение"
// @Success      200  {object}  response.DocsListResponse  "Получили документы"
// @Failure      400  {object}   response.APIResponse   "Неверный запрос"
// @Failure      401  {object}  response.APIResponse  "Не авторизовались"
// @Failure      500  {object}   response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/docs/{token}/{id} [get]
func (s *V1) GetDocsById(w http.ResponseWriter, r *http.Request) {
	
	q := r.URL.Query()

	token := q.Get("token")
	if token == "" {
		sendError(w, http.StatusUnauthorized, "Токен обязателен")
		return
	}

	_, err := s.U.GetSessionByToken(token)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Неверный токен")
		return
	}
	ID := q.Get("id")
	if ID == "" {
		sendError(w, http.StatusUnauthorized, "ID обязателен")
		return
	}
	var docs *response.Document
	// Лезем в кэш
	docs = s.C.GetDocsById(ID)
	// Не нашли лезем в базу
	if docs == nil {
		docs, err = s.F.GetDocById(ID)
		if err != nil {
			sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка получения документа: %s", err.Error()))
			return
		}
	}

	var resp any
	if !bytes.Equal(docs.JSONData, []byte(`{}`)) {
		resp = response.JSONResponse{
			Data: docs.JSONData,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		if len(docs.FileData) == 0 {
			sendError(w, http.StatusNotFound, "Файл отсутствует")
			return
		}
		mime := docs.Mime
		if mime == "" {
			mime = "application/octet-stream"
		}
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Length", strconv.Itoa(len(docs.FileData)))
		w.Header().Set(
			"Content-Disposition",
			fmt.Sprintf(`attachment; filename="%s"`, docs.Name),
		)
		_, _ = w.Write(docs.FileData)
	}

}

// DeleteDocs godoc
// @Summary      Удаление документа по ID
// @Description  Удаление документа по ID
// @Tags         docs
// @Accept 		 json
// @Produce      json
// @Param        token  query     string  true   "Токен авторизации"
// @Param        id  query     string     true  "ID документа"
// @Success      200  {file}  response.DynamicResponse "Удалили документ"
// @Failure      400  {object}   response.APIResponse   "Неверный запрос"
// @Failure      401  {object}  response.APIResponse  "Не авторизовались"
// @Failure      500  {object}   response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/docs/{token}/{id} [delete]
func (s *V1) DeleteDocs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	token := q.Get("token")
	if token == "" {
		sendError(w, http.StatusUnauthorized, "Токен обязателен")
		return
	}

	_, err := s.U.GetSessionByToken(token)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Неверный токен")
		return
	}
	ID := q.Get("id")
	if ID == "" {
		sendError(w, http.StatusUnauthorized, "ID обязателен")
		return
	}
	err = s.F.DeleteDocById(ID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка удаления документа: %s", err.Error()))
		return
	}
	// Удаляем из кэша
	s.C.DeleteData(ID)
	resp := response.NewDynamicResponse(token, true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func validateMeta(meta request.UploadMeta) error {
	if meta.Name == "" {
		return fmt.Errorf("поле 'name' обязательно")
	}

	if len(meta.Name) > 255 {
		return fmt.Errorf("имя файла слишком длинное (макс. 255 символов)")
	}
	// Здесь надо еще проверить реально ли это meme
	if meta.Mime == "" {
		return fmt.Errorf("поле 'mime' обязательно")
	}

	return nil
}
