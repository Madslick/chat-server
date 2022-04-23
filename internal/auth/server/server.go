package server

import (
	"github.com/Madslick/chit-chat-go/internal/auth/connectors"
	"github.com/Madslick/chit-chat-go/internal/auth/pkg"
)

func NewServer(authConnector connectors.AccountConnector) *AuthServer {
	return &AuthServer{authConnector: authConnector}
}

type AuthServer struct {
	pkg.UnimplementedAuthServer

	authConnector connectors.AccountConnector
}
