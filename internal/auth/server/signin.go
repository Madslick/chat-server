package server

import (
	"context"
	"log"

	"github.com/Madslick/chit-chat-go/internal/auth/pkg"
)

func (as *AuthServer) SignIn(ctx context.Context, in *pkg.SignInRequest) (*pkg.Account, error) {
	account, err := as.authConnector.SignIn(in)
	if err != nil {
		log.Fatalf("Error occurred while signing in %s", err)
	}

	return account, nil
}
