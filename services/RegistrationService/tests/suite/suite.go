package suite

import (
	"context"
	loadconfig "github.com/kv1sidisi/shapeup/pkg/config"
	regv1 "github.com/kv1sidisi/shapeup/pkg/proto/usercreatesvc/pb"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg       *config.Config
	RegClient regv1.UserCreationClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := &config.Config{}
	loadconfig.MustLoadByPath("../config/local.yaml", cfg)

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
		regv1.NewUserCreationClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(int(cfg.GRPC.Port)))
}
