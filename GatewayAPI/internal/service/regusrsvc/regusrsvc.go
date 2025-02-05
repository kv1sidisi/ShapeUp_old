package regusrsvc

import (
	pbusrcreatesvc "GatewayAPI/api/grpccl/pb/usrcreatesvc"
	"context"
	"log/slog"
)

type RegUsrSvc struct {
	log    *slog.Logger
	client pbusrcreatesvc.UserCreationClient
}

// New creates RegUsrSvc service.
func New(log *slog.Logger, client pbusrcreatesvc.UserCreationClient) *RegUsrSvc {
	return &RegUsrSvc{
		log:    log,
		client: client,
	}
}

// RegisterUser method invokes grpc client of RegisterService to register new user.
func (ru *RegUsrSvc) RegisterUser(email string, password string) (resp *pbusrcreatesvc.RegisterResponse, err error) {
	const op = "reusrsvc.RegisterUser"
	log := ru.log.With(slog.String("op", op))

	log.Info("registering user: ", email, password)
	resp, err = ru.client.Register(context.Background(), &pbusrcreatesvc.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Error("register error", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("registered account ", slog.Int64("userId", resp.GetUserId()))
	return resp, nil
}
