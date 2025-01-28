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

// SendingService struct represents the sending service and it is implementation of bottom layer of sending method of application.
type SendingService struct {
	log *slog.Logger
	cfg *config.Config
}

// New returns a new instance of SendingService service.
func New(log *slog.Logger, cfg *config.Config) *SendingService {
	return &SendingService{
		log: log,
		cfg: cfg}
}

// SMTPSendNewEmail sends email through SMTP with smtp package.
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

	log.Info("sending email through SMTP smtp package")
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

// smtpAddress merges host and port from config.
func smtpAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.SMTP.YDX.Host, strconv.Itoa(int(cfg.SMTP.YDX.Port)))
}

// GoGetSendNewEmail sends email through SMTP with GoGet package.
func (ss *SendingService) GoGetSendNewEmail(ctx context.Context, email string, message string) error {
	const op = "service.GoGetSendNewEmail"
	log := ss.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("sending email through SMTP GoGetPackage")
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
