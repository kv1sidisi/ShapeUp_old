package sendsvc

import (
	"context"
	"github.com/go-gomail/gomail"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/config"
	"log/slog"
)

// SendSvc sending service.
type SendSvc struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *SendSvc {
	return &SendSvc{
		log: log,
		cfg: cfg}
}

// GoGetSendNewEmail sends email through SMTP with GoGet package.
//
// Returns:
//   - Error if: Fails to send email through SMTP.
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
		log.Error("sending error: ", err)
		return errdefs.ErrSendEmail
	}
	log.Info("email sent through SMTP GoGetPackage", slog.String("message", message))

	return nil
}
