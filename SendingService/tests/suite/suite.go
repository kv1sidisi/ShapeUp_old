package suite

import (
	sendv1 "SendingService/api/pb"
	"SendingService/internal/config"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg       *config.Config
	RegClient sendv1.SendingClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect to grpc server: %v", err)
	}

	return ctx, &Suite{
		t,
		cfg,
		sendv1.NewSendingClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(int(cfg.GRPC.Port)))
}
