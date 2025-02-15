package grpcsrv

import (
	"context"
	"fmt"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/api/grpc/pb/jwtsvc"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	"google.golang.org/grpc"
	"log/slog"
)

// JWTSvc service for serverAPI.
type JWTSvc interface {
	GenerateToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error)
	ValidateToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error)
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

// GenerateToken is the GRPC server handler method. Generates new JWT token.
//
// Returns:
//
//   - A pointer to TokenResponse if succeeded.
//
//   - An error if: request is invalid. Error while generating token through service.
func (s *ServerAPI) GenerateToken(ctx context.Context,
	req *jwtsvc.GenerateTokenRequest,
) (*jwtsvc.GenerateTokenResponse, error) {
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

	token, err := s.jwtService.GenerateToken(ctx, req.GetUid(), req.GetOperation(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	log.Info("access token generated successfully: ", token)

	return &jwtsvc.GenerateTokenResponse{
		Token: token,
	}, nil
}

// ValidateToken is the GRPC server handler method. Validates JWT token.
//
// Returns:
//
//   - A pointer to ValidateTokenResponse if succeeded.
//
//   - An error if: Token is invalid. Request is invalid. Error while validating token through service.
func (s *ServerAPI) ValidateToken(ctx context.Context,
	req *jwtsvc.ValidateTokenRequest,
) (*jwtsvc.ValidateTokenResponse, error) {
	const op = "grpcsrv.ValidateAccessToken"
	log := s.log.With(slog.String("op", op))

	if len(req.GetToken()) == 0 {
		log.Error("invalid token in request")
		return nil, errdefs.InvalidToken
	}

	uid, operation, err := s.jwtService.ValidateToken(ctx, req.GetToken(), s.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("token validated successfully uid: %s, operation: %s", uid, operation))

	return &jwtsvc.ValidateTokenResponse{
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
