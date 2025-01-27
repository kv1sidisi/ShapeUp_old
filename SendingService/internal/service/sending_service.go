package service

import (
	"SendingService/internal/config"
	"context"
	"github.com/go-gomail/gomail"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
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

func (ss *SendingService) SMTPSendNewEmail(ctx context.Context, email string, message string) error {
	const op = "service.SendEmail"

	log := ss.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	auth := smtp.PlainAuth(
		"",
		ss.cfg.SMTP.MailRu.Username,
		ss.cfg.SMTP.MailRu.Password,
		ss.cfg.SMTP.MailRu.Host,
	)
	smtpAddress := smtpAddress(ss.cfg)
	log.Info("smtp sets up on: " + smtpAddress)

	log.Info("sending email through SMTP")
	err := smtp.SendMail(
		smtpAddress,
		auth,
		ss.cfg.SMTP.MailRu.Username,
		[]string{email},
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return nil
}

func smtpAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.SMTP.YDX.Host, strconv.Itoa(int(cfg.SMTP.YDX.Port)))
}

func (ss *SendingService) GoGetSendNewEmail(ctx context.Context, email string, message string) error {
	const op = "service.GoGetSendNewEmail"
	log := ss.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("sending email through SMTP")
	m := gomail.NewMessage()
	m.SetHeader("From", ss.cfg.SMTP.MailRu.Username)
	m.SetHeader("To", email)
	m.SetBody("Body", message)

	d := gomail.NewDialer(ss.cfg.SMTP.MailRu.Host, int(ss.cfg.SMTP.MailRu.Port), ss.cfg.SMTP.MailRu.Username, ss.cfg.SMTP.MailRu.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
