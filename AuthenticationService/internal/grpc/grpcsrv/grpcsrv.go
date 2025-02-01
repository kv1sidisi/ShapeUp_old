package grpcsrv

import (
	pbauthsvc "AuthenticationService/api/pb/authentication_service"
	"AuthenticationService/api/pb/jwt_service"
	"AuthenticationService/api/pb/sending_service"
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
		jwtClient jwt_service.JWTClient,
	) (userId int64, accessToken string, refreshToken string, err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	pbauthsvc.UnimplementedAuthServer
	auth          Auth
	cfg           *config.Config
	log           *slog.Logger
	sendingClient sending_service.SendingClient
	jwtClient     jwt_service.JWTClient
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, auth Auth, cfg *config.Config, log *slog.Logger, sendingClient sending_service.SendingClient, jwtClient jwt_service.JWTClient) {
	pbauthsvc.RegisterAuthServer(gRPC, &serverAPI{
		auth:          auth,
		cfg:           cfg,
		log:           log,
		sendingClient: sendingClient,
		jwtClient:     jwtClient,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *pbauthsvc.LoginRequest,
) (*pbauthsvc.LoginResponse, error) {
	op := "server.Login"

	log := s.log.With(slog.String("op", op))

	log.Info("logging user: ")
	userId, jwt, refresh, err := s.auth.LoginUser(ctx, req.GetUsername(), req.GetPassword(), s.jwtClient)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Info("logged successfully userId: ", userId)

	return &pbauthsvc.LoginResponse{
		UserId:       userId,
		JwtToken:     jwt,
		RefreshToken: refresh,
	}, nil
}
