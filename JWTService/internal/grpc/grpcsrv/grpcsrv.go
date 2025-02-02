package grpcsrv

import (
	pbjwtsvc "JWTService/api/pb/jwtsvc"
	"JWTService/internal/config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
)

type JWTSvc interface {
	GenerateAccessToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	ValidateAccessToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error)
	ValidateRefreshToken(ctx context.Context, refreshToken string, secretKey string) (uid int64, operation string, err error)
	GenerateLink(ctx context.Context, linkBase string, uid int64, operation string, secretKey string) (string, error)
}

type ServerAPI struct {
	pbjwtsvc.UnimplementedJWTServer
	jwtService JWTSvc
	cfg        *config.Config
	log        *slog.Logger
}

func RegisterServer(grpcServer *grpc.Server, jwtService JWTSvc, cfg *config.Config, log *slog.Logger) {
	pbjwtsvc.RegisterJWTServer(grpcServer, &ServerAPI{
		jwtService: jwtService,
		cfg:        cfg,
		log:        log,
	})
}

func (s *ServerAPI) GenerateAccessToken(ctx context.Context,
	req *pbjwtsvc.AccessTokenRequest,
) (*pbjwtsvc.AccessTokenResponse, error) {
	const op = "grpcsrv.GenerateAccessToken"

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

	return &pbjwtsvc.AccessTokenResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) GenerateRefreshToken(ctx context.Context,
	req *pbjwtsvc.RefreshTokenRequest,
) (*pbjwtsvc.RefreshTokenResponse, error) {
	const op = "grpcsrv.GenerateRefreshToken"

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

	return &pbjwtsvc.RefreshTokenResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) ValidateAccessToken(ctx context.Context,
	req *pbjwtsvc.ValidateAccessTokenRequest,
) (*pbjwtsvc.ValidateAccessTokenResponse, error) {
	const op = "grpcsrv.ValidateAccessToken"

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

	return &pbjwtsvc.ValidateAccessTokenResponse{
		Operation: operation,
		Uid:       uid,
	}, nil
}

func (s *ServerAPI) ValidateRefreshToken(ctx context.Context,
	req *pbjwtsvc.ValidateRefreshTokenRequest,
) (*pbjwtsvc.ValidateRefreshTokenResponse, error) {
	const op = "grpcsrv.ValidateRefreshToken"

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

	return &pbjwtsvc.ValidateRefreshTokenResponse{
		Operation: operation,
		Uid:       uid,
	}, nil
}

func (s *ServerAPI) GenerateLink(ctx context.Context,
	req *pbjwtsvc.GenerateLinkRequest,
) (*pbjwtsvc.GenerateLinkResponse, error) {
	const op = "grpcsrv.GenerateLink"

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

	return &pbjwtsvc.GenerateLinkResponse{
		Link: link,
	}, nil
}

func validateGenerateLinkRequest(req *pbjwtsvc.GenerateLinkRequest) error {
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
