package sendsvc

import (
	"SendingService/internal/config"
	"context"
	"github.com/go-gomail/gomail"
	"log/slog"
)

// SendSvc struct represents the sending service and it is implementation of bottom layer of sending method of application.
type SendSvc struct {
	log *slog.Logger
	cfg *config.Config
}

// New returns a new instance of SendSvc service.
func New(log *slog.Logger, cfg *config.Config) *SendSvc {
	return &SendSvc{
		log: log,
		cfg: cfg}
}

// GoGetSendNewEmail sends email through SMTP with GoGet package.
func (ss *SendSvc) GoGetSendNewEmail(ctx context.Context, email string, message string) error {
	const op = "service.GoGetSendNewEmail"
	log := ss.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	m := gomail.NewMessage()
	m.SetHeader("From", ss.cfg.SMTP.MailRu.Username)
	m.SetHeader("To", email)
	m.SetBody("Body", message)

	d := gomail.NewDialer(ss.cfg.SMTP.MailRu.Host, int(ss.cfg.SMTP.MailRu.Port), ss.cfg.SMTP.MailRu.Username, ss.cfg.SMTP.MailRu.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	log.Info("email sent through SMTP GoGetPackage", slog.String("message", message))

	return nil
}
