package regusrsvc

import (
	"context"
	pbusrcreatesvc "github.com/kv1sidisi/shapeup/pkg/proto/usercreatesvc/pb"
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

// RegisterUser method invokes GRPCClientConfig client of RegisterService to register new user.
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

	log.Info("registered account ", slog.Any("userId", resp.GetUid()))
	return resp, nil
}
