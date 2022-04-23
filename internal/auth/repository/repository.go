package repository

import (
	"github.com/Madslick/chit-chat-go/internal/auth/datastructs"
	"github.com/Madslick/chit-chat-go/internal/shared/db"
)

type Repository interface {
	init()
	SignUp(email string, password string, first string, last string, phone string) (string, error)
	SignIn(email string, password string) (*datastructs.Account, error)
	SearchAccounts(searchQuery string) ([]*datastructs.Account, error)
}

func NewRepository(connection db.DbConnection) Repository {
	repo := &mongoRepository{conn: connection}
	repo.init()
	return repo
}
