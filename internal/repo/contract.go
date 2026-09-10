package repo

import (
	"errors"
	"test_astral/internal/controller/DTO/response"
)

type UserData struct {
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
	UserID       string `db:"id"`
}
type FileRepo interface {
	CreateDocument(doc response.Document) error
	ListDocuments(limit int) ([]response.Document, error)
	GetDocById(id string) (*response.Document, error)
	DeleteDocById(id string) error
}
type UserRepo interface {
	CreateUser(login string, password string) (string, error)
	GetDataUser(login string) (*UserData, error)
	InsertSession(userID, token string) error
	GetSessionByToken(token string) (string, error)
	DeleteToken(token string) (string, error)
}

var (
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidCredentials = errors.New("неверный логин или пароль")
	ErrUserAlreadyExists  = errors.New("пользователь уже существует")
	ErrDatabase           = errors.New("ошибка базы данных")
	ErrSessionNotFound    = errors.New("ссесия отсутсвует")
)
