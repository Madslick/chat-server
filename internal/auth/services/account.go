package services

import (
	"github.com/Madslick/chit-chat-go/internal/auth/datastructs"
	"github.com/Madslick/chit-chat-go/internal/auth/repository"
)

type AccountService interface {
	SignIn(email string, password string) (*datastructs.Account, error)
	SignUp(email string, password string, first string, last string, phone string) (string, error)
	SearchAccounts(searchQuery string) ([]*datastructs.Account, error)
}

type accountService struct {
	repo repository.Repository
}

func NewAccountService(repo repository.Repository) AccountService {
	return &accountService{repo: repo}
}

func (as *accountService) SignIn(email string, password string) (*datastructs.Account, error) {
	return as.repo.SignIn(email, password)

}

func (as *accountService) SignUp(email string, password string, first string, last string, phone string) (string, error) {
	return as.repo.SignUp(email, password, first, last, phone)
}

func (as *accountService) SearchAccounts(searchQuery string) ([]*datastructs.Account, error) {
	return as.repo.SearchAccounts(searchQuery)
}
