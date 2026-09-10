package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"test_astral/internal/repo"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db         repo.UserRepo
	adminToken string
}

func NewUserService(db repo.UserRepo, token string) *UserService {

	return &UserService{db: db, adminToken: token}
}

func (u *UserService) CreateUser(login string, password string) (string, error) {
	fmt.Println("Создаем в usecase")
	hash, err := hashPassword(password)
	if err != nil {
		return "", fmt.Errorf("Ошибка хэшированя пароля", err.Error())
	}
	fmt.Println("Создали пароль")

	return u.db.CreateUser(login, hash)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func generateAuthToken(login string) string {

	data := fmt.Sprintf("%s-%d-%s", login, time.Now().UnixNano(), "secret-salt")
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:32]
}

func (u *UserService) Auth(login string, password string) (string, error) {
	user, err := u.db.GetDataUser(login)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Неверный логин или пароль")
	}
	token := generateAuthToken(login)
	err = u.db.InsertSession(user.UserID, token)
	if err != nil {
		return "", fmt.Errorf("Не удалось создать сессию: %s", err.Error())
	}
	return token, nil
}
func (u *UserService) DeleteToken(token string) (string, error) {
	return u.db.DeleteToken(token)
}

func (u *UserService) GetAdminToken() string {
	return u.adminToken
}
func (u *UserService) GetSessionByToken(token string) (string, error) {
	return u.db.GetSessionByToken(token)
}
