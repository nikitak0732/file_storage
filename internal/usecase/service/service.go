package service

import (
	"test_astral/config"
	"test_astral/internal/repo"
)

type UseCase struct {
	User  *UserService
	Cache *Cache
	File  *FileService
}

func New(fileRepo repo.FileRepo, userRepo repo.UserRepo, cache *Cache, cfg *config.Config) *UseCase {
	return &UseCase{
		Cache: cache,
		File:  NewFileService(fileRepo),
		User:  NewUserService(userRepo, cfg.Base.AdminToken),
	}
}
