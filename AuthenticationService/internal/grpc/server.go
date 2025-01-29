package grpc

import (
	authv1 "AuthenticationService/api/pb/authentication"
	"AuthenticationService/api/pb/sending"
	"AuthenticationService/internal/config"
	"context"
	"google.golang.org/grpc"
	"log/slog"
)

// Auth interface represents upper layer of authentication methods of application.
type Auth interface {
	LoginUser(
		ctx context.Context,
		username string,
		password string,
	) (userId int64, err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth          Auth
	cfg           *config.Config
	log           *slog.Logger
	sendingClient sending.SendingClient
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, auth Auth, cfg *config.Config, log *slog.Logger, sendingClient sending.SendingClient) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{
		auth:          auth,
		cfg:           cfg,
		log:           log,
		sendingClient: sendingClient,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	op := "server.Login"

	log := s.log.With(slog.String("op", op))

	log.Info("logging user: ")

	//TODO: validate request

	//TODO: invoke service

	//TODO: return response
}
