package service_register_user

import (
	regv1 "GatewayAPI/pkg/grpc_client/pb"
	"context"
	"fmt"
	"log/slog"
)

type RegisterUser struct {
	log    *slog.Logger
	client regv1.UserCreationClient
}

// New creates RegisterUser service.
func New(log *slog.Logger, client regv1.UserCreationClient) *RegisterUser {
	return &RegisterUser{
		log:    log,
		client: client,
	}
}

// RegisterUser method invokes grpc client of RegisterService to register new user.
func (ru *RegisterUser) RegisterUser(email string, password string) error {

	ru.log.Info("registering user: ", email, password)
	resp, err := ru.client.Register(context.Background(), &regv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		ru.log.Error("register error", slog.String("error", err.Error()))
		return err
	}

	fmt.Println(resp.UserId)

	ru.log.Info("registered account ", slog.Int64("userId", resp.UserId))
	return nil
}
