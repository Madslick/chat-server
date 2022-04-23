package server

import (
	"context"
	"log"

	"github.com/Madslick/chit-chat-go/internal/auth/pkg"
)

func (as *AuthServer) SearchAccounts(ctx context.Context, in *pkg.SearchAccountsRequest) (*pkg.SearchAccountsResponse, error) {
	accounts, err := as.authConnector.SearchAccounts(in)
	if err != nil {
		log.Fatalf("Error occurred while searching accounts %s", err)
	}

	return accounts, nil
}
