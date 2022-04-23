package server

import (
	"context"

	"github.com/Madslick/chit-chat-go/internal/auth/pkg"
)

func (as *AuthServer) SignUp(ctx context.Context, request *pkg.SignUpRequest) (*pkg.SignUpResponse, error) {
	signupResponse, _ := as.authConnector.SignUp(request)
	return signupResponse, nil
}
