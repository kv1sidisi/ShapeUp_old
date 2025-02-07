package grpccl

import (
	pbauthsvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/authsvc"
	pbusrcreatesvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/usrcreatesvc"
	"github.com/kv1sidisi/shapeup/services/gtwapi/cmd/grpccl/consts"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/config"
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

	clients.Cl[consts.UsrCreateSvc] = InitUsrCreteCl(log, cfg)
	clients.Cl[consts.AuthSvc] = InitAuthCl(log, cfg)

	return clients
}

func InitUsrCreteCl(log *slog.Logger, cfg *config.Config) *ClConn {
	usrCreateSvcConn := mustConnectToGRPC(cfg.GRPC.UserCreationServiceAddress)
	usrCreateSvcClient := pbusrcreatesvc.NewUserCreationClient(usrCreateSvcConn)
	log.Info("GRPC RegistrationService connected", slog.String("address", cfg.GRPC.UserCreationServiceAddress))
	return &ClConn{
		Client: usrCreateSvcClient,
		Conn:   usrCreateSvcConn,
	}
}

func InitAuthCl(log *slog.Logger, cfg *config.Config) *ClConn {
	authSvcConn := mustConnectToGRPC(cfg.GRPC.AuthenticationServiceAddress)
	authSvcClient := pbauthsvc.NewAuthClient(authSvcConn)
	log.Info("GRPC AuthService connected", slog.String("address", cfg.GRPC.AuthenticationServiceAddress))
	return &ClConn{
		Client: authSvcClient,
		Conn:   authSvcConn,
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
