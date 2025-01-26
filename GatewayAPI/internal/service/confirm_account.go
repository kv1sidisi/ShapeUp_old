package service

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
	ca.log.Info("Sending token for confirmation", slog.String("token", token))

	resp, err := ca.client.Confirm(context.Background(), &regv1.ConfirmRequest{Jwt: token})
	if err != nil {
		ca.log.Error("Confirm error", slog.String("error", err.Error()))
		return err
	}

	fmt.Println(resp.UserId)

	ca.log.Info("ConfirmAccount response: ", resp.UserId)
	ca.log.Info("ConfirmAccount success")
	return nil
}
