package service

import (
	"SendingService/internal/config"
	"context"
	"log/slog"
	"net/smtp"
)

type SendingService struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *SendingService {
	return &SendingService{
		log: log,
		cfg: cfg}
}

func (ss *SendingService) SendNewEmail(ctx context.Context, email string, message string) error {
	const op = "service.SendEmail"

	log := ss.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	auth := smtp.PlainAuth(
		"",
		ss.cfg.SMTP.Username,
		ss.cfg.SMTP.Password,
		ss.cfg.SMTP.Host,
	)
	log.Info("sending email through SMTP")
	err := smtp.SendMail(
		ss.cfg.SMTP.HostPort,
		auth,
		ss.cfg.SMTP.Username,
		[]string{email},
		[]byte(message),
	)
	if err != nil {
		log.Error("failed to send email", err)
		return err
	}
	return nil
}
