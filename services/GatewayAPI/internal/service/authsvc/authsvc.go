package authsvc

import (
	"context"
	pbauthsvc "github.com/kv1sidisi/shapeup/pkg/proto/authsvc/pb"
	"log/slog"
)

type AuthSvc struct {
	log    *slog.Logger
	client pbauthsvc.AuthClient
}

// New creates AuthSvc service.
func New(log *slog.Logger, client pbauthsvc.AuthClient) *AuthSvc {
	return &AuthSvc{
		log:    log,
		client: client}
}

// Login method invokes GRPCClientConfig client of AuthenticationService to log user in.
func (as *AuthSvc) Login(username string, password string) (*pbauthsvc.LoginResponse, error) {
	const op = "authsvc.Login"
	log := as.log.With(slog.String("op", op))

	log.Info("logging user in: ", username, password)
	resp, err := as.client.Login(context.Background(), &pbauthsvc.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error("login error", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("logged in account ", slog.Any("userId", resp.GetUid()))
	return resp, nil
}
