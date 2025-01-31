package grpc_server

import (
	jwtv1 "JWTService/api/pb"
	"JWTService/internal/config"
	"context"
	"google.golang.org/grpc"
	"log/slog"
)

type JWT interface {
	GenerateAccessToken(ctx context.Context, uid int64) (string, error)
	GenerateRefreshToken(ctx context.Context, uid int64) (string, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (uid int64, operation string, err error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (uid int64, operation string, err error)
	GenerateLink(ctx context.Context, linkBase string, uid int64, operation string) (string, error)
}

type ServerAPI struct {
	jwtv1.UnimplementedJWTServer
	jwtService JWT
	cfg        *config.Config
	log        *slog.Logger
}

func RegisterServer(grpcServer *grpc.Server, jwtService JWT, cfg *config.Config, log *slog.Logger) {
	jwtv1.RegisterJWTServer(grpcServer, &ServerAPI{
		jwtService: jwtService,
		cfg:        cfg,
		log:        log,
	})
}

func (s *ServerAPI) GenerateAccessToken(ctx context.Context,
	req *jwtv1.AccessTokenRequest,
) (*jwtv1.AccessTokenResponse, error) {
	panic("implement me")
}

func (s *ServerAPI) GenerateRefreshToken(ctx context.Context,
	req *jwtv1.RefreshTokenRequest,
) (*jwtv1.RefreshTokenResponse, error) {
	panic("implement me")
}

func (s *ServerAPI) ValidateAccessToken(ctx context.Context,
	req *jwtv1.ValidateAccessTokenRequest,
) (*jwtv1.ValidateAccessTokenResponse, error) {
	panic("implement me")
}

func (s *ServerAPI) ValidateRefreshToken(ctx context.Context,
	req *jwtv1.ValidateRefreshTokenRequest,
) (*jwtv1.ValidateRefreshTokenResponse, error) {
	panic("implement me")
}

func (s *ServerAPI) GenerateLink(ctx context.Context,
	req *jwtv1.GenerateLinkRequest,
) (*jwtv1.GenerateLinkResponse, error) {
	panic("implement me")
}
