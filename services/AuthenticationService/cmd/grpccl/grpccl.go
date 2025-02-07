package grpccl

import (
	pbjwtsvc "github.com/kv1sidisi/shapeup/services/authsvc/api/grpccl/pb/jwtsvc"
	pbsendsvc "github.com/kv1sidisi/shapeup/services/authsvc/api/grpccl/pb/sendsvc"
	"github.com/kv1sidisi/shapeup/services/authsvc/cmd/grpccl/consts"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

// GRPCClients struct contains map of GRPC clients and their connections to GRPC server.
type GRPCClients struct {
	log *slog.Logger
	Cl  map[string]*ClConn
}

type ClConn struct {
	Client interface{}
	Conn   *grpc.ClientConn
}

func New(log *slog.Logger, cfg *config.Config) *GRPCClients {
	clients := &GRPCClients{
		log: log,
		Cl:  make(map[string]*ClConn),
	}

	clients.Cl[consts.JWTSvc] = InitJWTCl(log, cfg)
	clients.Cl[consts.SendSvc] = InitSendCl(log, cfg)

	return clients
}

func InitSendCl(log *slog.Logger, cfg *config.Config) *ClConn {
	sendingServiceConn := mustConnectToGRPC(cfg.GRPCClient.SendingServiceAddress)
	sendingClient := pbsendsvc.NewSendingClient(sendingServiceConn)
	log.Info("GRPC SendingService connected", slog.String("address", cfg.GRPCClient.SendingServiceAddress))
	return &ClConn{
		Client: sendingClient,
		Conn:   sendingServiceConn,
	}
}

func InitJWTCl(log *slog.Logger, cfg *config.Config) *ClConn {
	jwtServiceConn := mustConnectToGRPC(cfg.GRPCClient.JWTServiceAddress)
	jwtClient := pbjwtsvc.NewJWTClient(jwtServiceConn)
	log.Info("GRPC JWTService connected", slog.String("address", cfg.GRPCClient.JWTServiceAddress))
	return &ClConn{
		Client: jwtClient,
		Conn:   jwtServiceConn,
	}
}

func mustConnectToGRPC(address string) *grpc.ClientConn {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *GRPCClients) Close() {
	for name, clientConn := range c.Cl {
		if err := clientConn.Conn.Close(); err != nil {
			c.log.Error("failed to close connection for service", slog.String("service", name), err)
		}
	}
}
