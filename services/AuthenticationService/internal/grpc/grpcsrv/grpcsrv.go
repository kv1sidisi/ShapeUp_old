package grpcsrv

import (
	"context"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/pkg/proto/authsvc/pb"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/config"
	"google.golang.org/grpc"
	"log/slog"
)

// AuthSvc service for serverAPI.
type AuthSvc interface {
	LoginUser(
		ctx context.Context,
		username string,
		password string,
	) (userId int64, accessToken string, refreshToken string, err error)
}

// serverAPI handler for the gRPC server.
type serverAPI struct {
	authsvc.UnimplementedAuthServer
	auth AuthSvc
	cfg  *config.Config
	log  *slog.Logger
}

// RegisterServer registers services in the GRPC server.
//
// Returns serverAPI as handler for GRPC server.
func RegisterServer(gRPC *grpc.Server, auth AuthSvc, cfg *config.Config, log *slog.Logger) {
	authsvc.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
		cfg:  cfg,
		log:  log,
	})
}

// Login is the GRPC server handler method. Logs user in.
//
// Returns:
//
//   - A pointer to LoginResponse if successful.
//
//   - An error if: Request is invalid. Error while logging user in through service.
func (s *serverAPI) Login(
	ctx context.Context,
	req *authsvc.LoginRequest,
) (*authsvc.LoginResponse, error) {
	const op = "grpcsrv.Login"

	log := s.log.With(slog.String("op", op))

	if len(req.Username) == 0 || len(req.Password) == 0 {
		log.Info("got wrong credentials: ", req.Username, req.Password)
		return nil, errdefs.InvalidCredentials
	}

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
