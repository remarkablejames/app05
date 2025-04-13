package storage

import (
	"app05/internal/core/domain/repositories"
	"app05/internal/infrastructure/storage/postgres/repo_impl"
	"database/sql"
)

type Storage struct {
	User    repositories.UserRepository
	Session repositories.SessionRepository
	Post    repositories.PostRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User:    repo_impl.NewUserRepository(db),
		Session: repo_impl.NewSessionRepository(db),
		Post:    repo_impl.NewPostRepository(db),
	}
}
