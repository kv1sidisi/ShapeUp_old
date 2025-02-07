package grpcsrv

import (
	"context"
	"github.com/kv1sidisi/shapeup/services/authsvc/api/grpc/pb/authsvc"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/config"
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
	authsvc.UnimplementedAuthServer
	auth AuthSvc
	cfg  *config.Config
	log  *slog.Logger
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, auth AuthSvc, cfg *config.Config, log *slog.Logger) {
	authsvc.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
		cfg:  cfg,
		log:  log,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authsvc.LoginRequest,
) (*authsvc.LoginResponse, error) {
	const op = "grpcsrv.Login"

	log := s.log.With(slog.String("op", op))

	uid, jwt, refresh, err := s.auth.LoginUser(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	log.Info("logged successfully userId: ", uid)

	return &authsvc.LoginResponse{
		UserId:       uid,
		JwtToken:     jwt,
		RefreshToken: refresh,
	}, nil
}
