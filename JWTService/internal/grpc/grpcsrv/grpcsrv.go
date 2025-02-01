package grpcsrv

import (
	pbjwt "JWTService/api/pb"
	"JWTService/internal/config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
)

type JWT interface {
	GenerateAccessToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	ValidateAccessToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error)
	ValidateRefreshToken(ctx context.Context, refreshToken string, secretKey string) (uid int64, operation string, err error)
	GenerateLink(ctx context.Context, linkBase string, uid int64, operation string, secretKey string) (string, error)
}

type ServerAPI struct {
	pbjwt.UnimplementedJWTServer
	jwtService JWT
	cfg        *config.Config
	log        *slog.Logger
}

func RegisterServer(grpcServer *grpc.Server, jwtService JWT, cfg *config.Config, log *slog.Logger) {
	pbjwt.RegisterJWTServer(grpcServer, &ServerAPI{
		jwtService: jwtService,
		cfg:        cfg,
		log:        log,
	})
}

func (s *ServerAPI) GenerateAccessToken(ctx context.Context,
	req *pbjwt.AccessTokenRequest,
) (*pbjwt.AccessTokenResponse, error) {
	const op = "server.GenerateAccessToken"

	log := s.log.With(slog.String("op", op))

	if req.GetUid() == 0 {
		log.Error("invalid uid in request:", req.GetUid())
		return nil, fmt.Errorf("invalid uid in request")
	}
	if len(req.GetOperation()) == 0 {
		log.Error("invalid operation in request: ", req.GetUid())
		return nil, fmt.Errorf("invalid operation in request")
	}

	token, err := s.jwtService.GenerateAccessToken(ctx, req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("access token generated successfully: ", token)

	return &pbjwt.AccessTokenResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) GenerateRefreshToken(ctx context.Context,
	req *pbjwt.RefreshTokenRequest,
) (*pbjwt.RefreshTokenResponse, error) {
	const op = "server.GenerateRefreshToken"

	log := s.log.With(slog.String("op", op))

	if req.GetUid() == 0 {
		log.Error("invalid uid in request")
		return nil, fmt.Errorf("invalid uid in request")
	}
	if len(req.GetOperation()) == 0 {
		log.Error("invalid operation in request")
		return nil, fmt.Errorf("invalid operation in request")
	}

	token, err := s.jwtService.GenerateRefreshToken(ctx, req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("refresh token generated successfully: ", token)

	return &pbjwt.RefreshTokenResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) ValidateAccessToken(ctx context.Context,
	req *pbjwt.ValidateAccessTokenRequest,
) (*pbjwt.ValidateAccessTokenResponse, error) {
	const op = "server.ValidateAccessToken"

	log := s.log.With(slog.String("op", op))

	if len(req.GetToken()) == 0 {
		log.Error("invalid token in request")
	}

	uid, operation, err := s.jwtService.ValidateAccessToken(ctx, req.GetToken(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info(fmt.Sprintf("token validated successfully uid: %s, operation: %s", uid, operation))

	return &pbjwt.ValidateAccessTokenResponse{
		Operation: operation,
		Uid:       uid,
	}, nil
}

func (s *ServerAPI) ValidateRefreshToken(ctx context.Context,
	req *pbjwt.ValidateRefreshTokenRequest,
) (*pbjwt.ValidateRefreshTokenResponse, error) {
	const op = "server.ValidateRefreshToken"

	log := s.log.With(slog.String("op", op))

	if len(req.GetToken()) == 0 {
		log.Error("invalid token in request")
	}

	uid, operation, err := s.jwtService.ValidateRefreshToken(ctx, req.GetToken(), s.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info(fmt.Sprintf("token validated successfully uid: %s, operation: %s", uid, operation))

	return &pbjwt.ValidateRefreshTokenResponse{
		Operation: operation,
		Uid:       uid,
	}, nil
}

func (s *ServerAPI) GenerateLink(ctx context.Context,
	req *pbjwt.GenerateLinkRequest,
) (*pbjwt.GenerateLinkResponse, error) {
	const op = "server.GenerateLink"

	log := s.log.With(slog.String("op", op))

	if err := validateGenerateLinkRequest(req); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	link, err := s.jwtService.GenerateLink(ctx, req.GetLinkBase(), req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("link generated successfully: ", link)

	return &pbjwt.GenerateLinkResponse{
		Link: link,
	}, nil
}

func validateGenerateLinkRequest(req *pbjwt.GenerateLinkRequest) error {
	if req.GetUid() == 0 {
		return fmt.Errorf("invalid uid in request")
	}

	if req.GetOperation() == "" {
		return fmt.Errorf("invalid operation in request")
	}

	if req.GetLinkBase() == "" {
		return fmt.Errorf("invalid linkBase in request")
	}

	return nil
}
