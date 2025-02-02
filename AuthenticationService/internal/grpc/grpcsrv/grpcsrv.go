package grpcsrv

import (
	pbauthsvc "AuthenticationService/api/pb/authsvc"
	"AuthenticationService/internal/config"
	"context"
	"google.golang.org/grpc"
	"log/slog"
)

// AuthSvc interface represents upper layer of authentication methods of application.
type AuthSvc interface {
	LoginUser(
		ctx context.Context,
		username string,
		password string,
	) (userId int64, accessToken string, refreshToken string, err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	pbauthsvc.UnimplementedAuthServer
	auth AuthSvc
	cfg  *config.Config
	log  *slog.Logger
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, auth AuthSvc, cfg *config.Config, log *slog.Logger) {
	pbauthsvc.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
		cfg:  cfg,
		log:  log,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *pbauthsvc.LoginRequest,
) (*pbauthsvc.LoginResponse, error) {
	const op = "grpcsrv.Login"

	log := s.log.With(slog.String("op", op))

	uid, jwt, refresh, err := s.auth.LoginUser(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Info("logged successfully userId: ", uid)

	return &pbauthsvc.LoginResponse{
		UserId:       uid,
		JwtToken:     jwt,
		RefreshToken: refresh,
	}, nil
}
