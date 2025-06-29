package confaccsvc

import (
	"context"
	pbusrcreatesvc "github.com/kv1sidisi/shapeup/pkg/proto/usercreatesvc/pb"
	"log/slog"
)

type ConfAccSvc struct {
	log    *slog.Logger
	client pbusrcreatesvc.UserCreationClient
}

// New creates ConfAccSvc service
func New(log *slog.Logger, client pbusrcreatesvc.UserCreationClient) *ConfAccSvc {
	return &ConfAccSvc{
		log:    log,
		client: client,
	}
}

// ConfAccSvc method invokes grpc client of RegistrationService to confirm account
func (ca *ConfAccSvc) ConfirmAccount(token string) error {
	const op = "confaccsvc.ConfirmAccount"
	log := ca.log.With(slog.String("op", op))

	resp, err := ca.client.Confirm(context.Background(), &pbusrcreatesvc.ConfirmRequest{Jwt: token})
	if err != nil {
		log.Error("confirm error", slog.String("error", err.Error()))
		return err
	}

	log.Info("confirmed account ", slog.Any("userId", resp.GetUid()))
	return nil
}
