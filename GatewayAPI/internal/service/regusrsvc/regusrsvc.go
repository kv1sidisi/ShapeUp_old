package regusrsvc

import (
	pbusrcreatesvc "GatewayAPI/api/grpccl/pb/usrcreatesvc"
	"context"
	"fmt"
	"log/slog"
)

type RegisterUser struct {
	log    *slog.Logger
	client pbusrcreatesvc.UserCreationClient
}

// New creates RegisterUser service.
func New(log *slog.Logger, client pbusrcreatesvc.UserCreationClient) *RegisterUser {
	return &RegisterUser{
		log:    log,
		client: client,
	}
}

// RegisterUser method invokes grpc client of RegisterService to register new user.
func (ru *RegisterUser) RegisterUser(email string, password string) (resp *pbusrcreatesvc.RegisterResponse, err error) {

	ru.log.Info("registering user: ", email, password)
	resp, err = ru.client.Register(context.Background(), &pbusrcreatesvc.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		ru.log.Error("register error", slog.String("error", err.Error()))
		return nil, err
	}

	fmt.Println(resp.UserId)

	ru.log.Info("registered account ", slog.Int64("userId", resp.UserId))
	return resp, nil
}
