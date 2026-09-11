package api

import (
	"encoding/json"
	"net/http"
	"test_astral/internal/controller/DTO/response"
	"test_astral/internal/usecase"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type V1 struct {
	U usecase.User
	F usecase.File
	C usecase.Cache
}

func NewRoutes(f usecase.File, u usecase.User, c usecase.Cache) *V1 {
	r := &V1{F: f, U: u, C: c}
	return r
}

func (s *V1) RegisterRoutes(mux *http.ServeMux) {

	r := chi.NewRouter()
	s.RegisterRoutesFile(r)
	s.RegisterRoutesUser(r)
	r.Get("/check", s.HealthHandler)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Ручная отдача swagger.json
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/docs/swagger.json")
	})

	mux.Handle("/", r)

}

func (s *V1) HealthHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "processed",
		"data":   req,
	})
}

func sendError(w http.ResponseWriter, code int, text string) {
	resp := &response.APIResponse{
		Error: &response.ErrorDetail{
			Code: code,
			Text: text,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// func sendData(w http.ResponseWriter, code int, text string) {
// 	resp := &response.APIResponse{
// 		Error: &response.ErrorDetail{
// 			Code: code,
// 			Text: text,
// 		},
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(resp)
// }
