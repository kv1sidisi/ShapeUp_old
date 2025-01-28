package service_confirm_account

import (
	regv1 "GatewayAPI/pkg/grpc_client/pb"
	"context"
	"fmt"
	"log/slog"
)

type ConfirmAccount struct {
	log    *slog.Logger
	client regv1.UserCreationClient
}

func New(log *slog.Logger, client regv1.UserCreationClient) *ConfirmAccount {
	return &ConfirmAccount{
		log:    log,
		client: client,
	}
}

func (ca *ConfirmAccount) ConfirmAccount(token string) error {
	ca.log.Info("sending token for confirmation", slog.String("token", token))

	resp, err := ca.client.Confirm(context.Background(), &regv1.ConfirmRequest{Jwt: token})
	if err != nil {
		ca.log.Error("confirm error", slog.String("error", err.Error()))
		return err
	}

	fmt.Println(resp.UserId)

	ca.log.Info("confirmed account ", slog.Int64("userId", resp.UserId))
	return nil
}
