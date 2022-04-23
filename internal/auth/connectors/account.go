package connectors

import (
	"log"

	"github.com/Madslick/chit-chat-go/internal/auth/pkg"
	"github.com/Madslick/chit-chat-go/internal/auth/services"
)

type AccountConnector interface {
	SignUp(request *pkg.SignUpRequest) (*pkg.SignUpResponse, error)
	SignIn(request *pkg.SignInRequest) (*pkg.Account, error)
	SearchAccounts(request *pkg.SearchAccountsRequest) (*pkg.SearchAccountsResponse, error)
}

type accountConnector struct {
	accountService services.AccountService
}

func NewAccountConnector(accService services.AccountService) AccountConnector {
	return &accountConnector{
		accountService: accService,
	}
}

func (ac *accountConnector) SignUp(request *pkg.SignUpRequest) (*pkg.SignUpResponse, error) {
	accountId, err := ac.accountService.SignUp(request.GetEmail(), request.GetPassword(), request.GetFirstName(), request.GetPassword(), request.GetPhoneNumber())
	if err != nil {
		log.Fatalf("Unable to sign up %s: %s", request.GetEmail(), err)
		return nil, err
	}

	signupResponse := pkg.SignUpResponse{Id: accountId}

	return &signupResponse, nil
}

func (ac *accountConnector) SignIn(request *pkg.SignInRequest) (*pkg.Account, error) {
	account, err := ac.accountService.SignIn(request.GetEmail(), request.GetPassword())
	if err != nil {
		log.Fatalf("Unable to log in %s: %s", request.GetEmail(), err)
		return nil, err
	}

	Account := pkg.Account{
		Id:          account.Id,
		Email:       account.Email,
		FirstName:   account.First,
		LastName:    account.Last,
		PhoneNumber: account.Phone,
	}

	return &Account, nil
}

func (ac *accountConnector) SearchAccounts(request *pkg.SearchAccountsRequest) (*pkg.SearchAccountsResponse, error) {

	_, err := ac.accountService.SearchAccounts(request.GetSearchQuery())

	return &pkg.SearchAccountsResponse{}, err
}
