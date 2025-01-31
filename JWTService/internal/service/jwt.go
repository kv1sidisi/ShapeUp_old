package service

import (
	"JWTService/internal/config"
	"context"
	"log/slog"
)

type JWTService struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *JWTService {
	return &JWTService{log: log, cfg: cfg}
}

func (s *JWTService) GenerateAccessToken(ctx context.Context, uid int64) (string, error) {
	panic("implement me")
}

func (s *JWTService) GenerateRefreshToken(ctx context.Context, uid int64) (string, error) {
	panic("implement me")
}

func (s *JWTService) ValidateAccessToken(ctx context.Context, accessToken string) (uid int64, operation string, err error) {
	panic("implement me")
}

func (s *JWTService) ValidateRefreshToken(ctx context.Context, refreshToken string) (uid int64, operation string, err error) {
	panic("implement me")
}

func (s *JWTService) GenerateLink(ctx context.Context, linkBase string, uid int64, operation string) (string, error) {
	panic("implement me")
}
