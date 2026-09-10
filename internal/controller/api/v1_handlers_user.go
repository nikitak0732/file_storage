package api

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unicode"

	"net/http"
	"test_astral/internal/controller/DTO/request"
	"test_astral/internal/controller/DTO/response"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

func (s *V1) RegisterRoutesUser(r chi.Router) {
	r.Post("/api/register", s.Register)
	r.Post("/api/auth", s.Auth)
	r.Delete("/api/auth/{token}", s.DeleteToken)
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя и возвращает логин
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body request.Registration true "Данные нового пользака"
// @Success      200  {object}  response.APIResponse  "Зарегестрировались"
// @Failure      400  {object}   response.APIResponse   "Неверный запрос"
// @Failure      500  {object}   response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/register [post]
func (s *V1) Register(w http.ResponseWriter, r *http.Request) {

	var req request.Registration
	// Проверим структуру
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Ошибка создания запроса. Переданная структура не соответствует ожидаемой: %s", err.Error()))
		return
	}
	defer r.Body.Close()
	// Проверим токен
	fmt.Println(req.Token, s.U.GetAdminToken())
	if req.Token != s.U.GetAdminToken() {
		sendError(w, http.StatusBadRequest, "Неверный токен администратора")
		return
	}
	if !ValidateLogin(req.Login) {
		sendError(w, http.StatusBadRequest,
			"Логин должен содержать минимум 8 символов (только латиница и цифры)")
		return
	}
	if !ValidatePassword(req.Password) {
		sendError(w, http.StatusBadRequest,
			"Пароль должен содержать минимум 8 символов, минимум 2 буквы в разных регистрах, минимум 1 цифру и минимум 1 спецсимвол")
		return
	}
	fmt.Println("Создаем пользователя...")
	login, err := s.U.CreateUser(req.Login, req.Password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка сохранения пользователя: %s", err.Error()))
		return
	}
	resp := &response.APIResponse{
		Response: &response.ResponseDetail{
			Login: login,
		},
	}
	fmt.Println("Создали пользователя")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteToken godoc
// @Summary      Удаление токена
// @Description  Авторизация пользователя и возврат токена
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        token query string true "Токен"
// @Success      200  {object}  response.APIResponse  "Удалили"
// @Failure      400  {object}   response.APIResponse   "Неверный запрос"
// @Failure      500  {object}   response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/auth/{token} [delete]
func (s *V1) DeleteToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		sendError(w, http.StatusBadRequest, "Ошибка создания запроса: не передан токен")
		return
	}

	// Здесь по хорошему еще надо проверить на то что переданная стринга это теоретически токен что бы здесь 400 вернуть

	token, err := s.U.DeleteToken(token)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Внутрениияя ошибка сервера при удалении токена: %s", err.Error()))
		return
	}
	resp := response.NewDynamicResponse(token, true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Auth godoc
// @Summary      Авторизация пользователя
// @Description  Авторизация пользователя и возврат токена
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body request.Login true "Данные на авторизацию"
// @Success      200  {object}  response.APIResponse  "Авторизовались"
// @Failure      400  {object}   response.APIResponse   "Неверный запрос"
// @Failure      500  {object}   response.APIResponse  "Внутренняя ошибка сервера"
// @Router       /api/auth [post]
func (s *V1) Auth(w http.ResponseWriter, r *http.Request) {
	var req request.Login
	// Проверим структуру
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Ошибка создания запроса. Переданная структура не соответствует ожидаемой: %s", err.Error()))
		return
	}
	if req.Login == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Логин и пароль обязательны")
		return
	}
	token, err := s.U.Auth(req.Login, req.Password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка авторизации пользователя: %s", err.Error()))
		return
	}
	resp := &response.APIResponse{
		Response: &response.ResponseDetail{
			Token: token,
		},
	}
	fmt.Println("Авторизовались")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ValidateLogin(login string) bool {
	if len(login) < 8 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", login)
	return matched
}

func ValidatePassword(password string) bool {

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
