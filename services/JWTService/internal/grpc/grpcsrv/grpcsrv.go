package grpcsrv

import (
	"context"
	"fmt"
	"github.com/kv1sidisi/shapeup/libs/common/errdefs"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/api/grpc/pb/jwtsvc"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	"google.golang.org/grpc"
	"log/slog"
)

// JWTSvc service for serverAPI.
type JWTSvc interface {
	GenerateAccessToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	ValidateAccessToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error)
	ValidateRefreshToken(ctx context.Context, refreshToken string, secretKey string) (uid int64, operation string, err error)
	GenerateLink(ctx context.Context, linkBase string, uid int64, operation string, secretKey string) (string, error)
}

// serverAPI handler for the gRPC server.
type ServerAPI struct {
	jwtsvc.UnimplementedJWTServer
	jwtService JWTSvc
	cfg        *config.Config
	log        *slog.Logger
}

// RegisterServer registers services in the GRPC server.
//
// Returns serverAPI as handler for GRPC server.
func RegisterServer(grpcServer *grpc.Server, jwtService JWTSvc, cfg *config.Config, log *slog.Logger) {
	jwtsvc.RegisterJWTServer(grpcServer, &ServerAPI{
		jwtService: jwtService,
		cfg:        cfg,
		log:        log,
	})
}

// GenerateAccessToken is the GRPC server handler method. Generates new access JWT token.
//
// Returns:
//
//   - A pointer to AccessTokenResponse if succeeded.
//
//   - An error if: request is invalid. Error while generating access token through service.
func (s *ServerAPI) GenerateAccessToken(ctx context.Context,
	req *jwtsvc.AccessTokenRequest,
) (*jwtsvc.AccessTokenResponse, error) {
	const op = "grpcsrv.GenerateAccessToken"
	log := s.log.With(slog.String("op", op))

	if req.GetUid() == 0 {
		log.Error("invalid uid in request:", req.GetUid())
		return nil, errdefs.InvalidUserId
	}
	if len(req.GetOperation()) == 0 {
		log.Error("invalid operation in request: ", req.GetUid())
		return nil, errdefs.InvalidOperationType
	}

	token, err := s.jwtService.GenerateAccessToken(ctx, req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	log.Info("access token generated successfully: ", token)

	return &jwtsvc.AccessTokenResponse{
		Token: token,
	}, nil
}

// GenerateRefreshToken is the GRPC server handler method. Generates new refresh JWT token.
//
// Returns:
//
//   - A pointer to RefreshTokenResponse if succeeded.
//
//   - An error if: request is invalid. Error while generating refresh token through service.
func (s *ServerAPI) GenerateRefreshToken(ctx context.Context,
	req *jwtsvc.RefreshTokenRequest,
) (*jwtsvc.RefreshTokenResponse, error) {
	const op = "grpcsrv.GenerateRefreshToken"

	log := s.log.With(slog.String("op", op))

	if req.GetUid() == 0 {
		log.Error("invalid uid in request:", req.GetUid())
		return nil, errdefs.InvalidUserId
	}
	if len(req.GetOperation()) == 0 {
		log.Error("invalid operation in request: ", req.GetUid())
		return nil, errdefs.InvalidOperationType
	}

	token, err := s.jwtService.GenerateRefreshToken(ctx, req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	log.Info("refresh token generated successfully: ", token)

	return &jwtsvc.RefreshTokenResponse{
		Token: token,
	}, nil
}

// ValidateAccessToken is the GRPC server handler method. Validates access JWT token.
//
// Returns:
//
//   - A pointer to ValidateAccessTokenResponse if succeeded.
//
//   - An error if: Access token is invalid. Request is invalid. Error while validating access token through service.
func (s *ServerAPI) ValidateAccessToken(ctx context.Context,
	req *jwtsvc.ValidateAccessTokenRequest,
) (*jwtsvc.ValidateAccessTokenResponse, error) {
	const op = "grpcsrv.ValidateAccessToken"
	log := s.log.With(slog.String("op", op))

	if len(req.GetToken()) == 0 {
		log.Error("invalid token in request")
		return nil, errdefs.InvalidToken
	}

	uid, operation, err := s.jwtService.ValidateAccessToken(ctx, req.GetToken(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("token validated successfully uid: %s, operation: %s", uid, operation))

	return &jwtsvc.ValidateAccessTokenResponse{
		Operation: operation,
		Uid:       uid,
	}, nil
}

// ValidateRefreshToken is the GRPC server handler method. Validates refresh JWT token.
//
// Returns:
//
//   - A pointer to ValidateRefreshTokenResponse if succeeded.
//
//   - An error if: Refresh token is invalid. Request is invalid. Error while validating refresh token through service.
func (s *ServerAPI) ValidateRefreshToken(ctx context.Context,
	req *jwtsvc.ValidateRefreshTokenRequest,
) (*jwtsvc.ValidateRefreshTokenResponse, error) {
	const op = "grpcsrv.ValidateRefreshToken"
	log := s.log.With(slog.String("op", op))

	if len(req.GetToken()) == 0 {
		log.Error("invalid token in request")
		return nil, errdefs.InvalidToken
	}

	uid, operation, err := s.jwtService.ValidateRefreshToken(ctx, req.GetToken(), s.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("token validated successfully uid: %s, operation: %s", uid, operation))

	return &jwtsvc.ValidateRefreshTokenResponse{
		Operation: operation,
		Uid:       uid,
	}, nil
}

// GenerateLink is the GRPC server handler method. Generates link from userId and operation type.
//
// Returns:
//
//   - A pointer to GenerateLinkResponse if succeeded.
//
//   - An error if: Request is invalid. Error while generating link through service.
func (s *ServerAPI) GenerateLink(ctx context.Context,
	req *jwtsvc.GenerateLinkRequest,
) (*jwtsvc.GenerateLinkResponse, error) {
	const op = "grpcsrv.GenerateLink"
	log := s.log.With(slog.String("op", op))

	if err := validateGenerateLinkRequest(req); err != nil {
		log.Error("failed to validate generate link request")
		return nil, err
	}

	link, err := s.jwtService.GenerateLink(ctx, req.GetLinkBase(), req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)

	if err != nil {
		return nil, err
	}

	log.Info("link generated successfully: ", link)

	return &jwtsvc.GenerateLinkResponse{
		Link: link,
	}, nil
}

func validateGenerateLinkRequest(req *jwtsvc.GenerateLinkRequest) error {
	if req.GetUid() == 0 {
		return errdefs.InvalidUserId
	}

	if req.GetOperation() == "" {
		return errdefs.InvalidOperationType
	}

	if req.GetLinkBase() == "" {
		return errdefs.InvalidLinkBase
	}

	return nil
}
