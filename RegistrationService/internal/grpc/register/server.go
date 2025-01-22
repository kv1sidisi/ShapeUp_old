package register

import (
	regv1 "RegistrationService/api/pb"
	"context"
	"google.golang.org/grpc"
)

type serverAPI struct {
	regv1.UnimplementedRegistrationServer
}

func Register(gRPC *grpc.Server) {
	regv1.RegisterRegistrationServer(gRPC, &serverAPI{})
}
func (s *serverAPI) Register(
	ctx context.Context,
	req *regv1.RegisterRequest,
) (*regv1.RegisterResponse, error) {
	panic("implement me")
}
