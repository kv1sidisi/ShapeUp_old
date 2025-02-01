package confaccsvc

import (
	pbgtway "GatewayAPI/api/grpc_client/pb"
	"context"
	"fmt"
	"log/slog"
)

type ConfirmAccount struct {
	log    *slog.Logger
	client pbgtway.UserCreationClient
}

// New creates ConfirmAccount service
func New(log *slog.Logger, client pbgtway.UserCreationClient) *ConfirmAccount {
	return &ConfirmAccount{
		log:    log,
		client: client,
	}
}

// ConfirmAccount method invokes grpc client of RegistrationService to confirm account
func (ca *ConfirmAccount) ConfirmAccount(token string) error {
	ca.log.Info("sending token for confirmation", slog.String("token", token))

	resp, err := ca.client.Confirm(context.Background(), &pbgtway.ConfirmRequest{Jwt: token})
	if err != nil {
		ca.log.Error("confirm error", slog.String("error", err.Error()))
		return err
	}

	fmt.Println(resp.UserId)

	ca.log.Info("confirmed account ", slog.Int64("userId", resp.UserId))
	return nil
}
