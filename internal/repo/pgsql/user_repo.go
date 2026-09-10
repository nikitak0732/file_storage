package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"test_astral/internal/repo"

	"test_astral/pkg/postgres"
)

type UserRepo struct {
	pg *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {

	return &UserRepo{pg}
}

func (u *UserRepo) CreateUser(login string, password string) (string, error) {
	fmt.Println("База")
	_, err := u.GetDataUser(login)
	fmt.Println(err)
	if !errors.Is(err, repo.ErrUserNotFound) {
		return "", fmt.Errorf("Юзер с таким логином уже существует: %s", err.Error())
	}

	_, err = u.pg.Postgres.Exec("INSERT INTO users (login, password_hash) VALUES ($1, $2)", login, password)
	if err != nil {
		return "", fmt.Errorf("Не удалось добавить пользователя: %s", err.Error())
	}
	return login, nil
}

func (u *UserRepo) GetDataUser(login string) (*repo.UserData, error) {
	var user repo.UserData
	err := u.pg.Postgres.QueryRow("select password_hash, id from users where login = $1", login).Scan(&user.PasswordHash, &user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, repo.ErrUserNotFound
		}

		return nil, repo.ErrDatabase
	}
	fmt.Println("Нашли юзерка")
	return &repo.UserData{
		Login:        login,
		PasswordHash: user.PasswordHash,
		UserID:       user.UserID,
	}, nil
}
func (u *UserRepo) InsertSession(userID, token string) error {
	_, err := u.pg.Postgres.Exec("INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)", userID, token)
	if err != nil {
		return fmt.Errorf("Не удалось создать сессию: %s", err.Error())
	}
	return nil
}

func (u *UserRepo) GetSessionByToken(token string) (string, error) {
	id := ""
	err := u.pg.Postgres.QueryRow(`
        SELECT user_id
        FROM sessions 
        WHERE token_hash = $1
    `, token).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repo.ErrSessionNotFound
		}
		return "", repo.ErrDatabase
	}

	return id, nil
}

func (u *UserRepo) DeleteToken(token string) (string, error) {
	query := `delete from sessions where token_hash = $1`
	_, err := u.pg.Postgres.Exec(query, token)
	if err != nil {
		return "", fmt.Errorf("Не удалось удалить сессию %s", err.Error())
	}
	return token, nil

}
