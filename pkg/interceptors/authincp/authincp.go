package authincp

import (
	"context"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	pbjwtsvc "github.com/kv1sidisi/shapeup/pkg/proto/jwtsvc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

func AuthInterceptor(log *slog.Logger, jwtsvc pbjwtsvc.JWTClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		const op = "authincp"
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Error("missing metadata for ", slog.String("interceptor", op))
			return nil, errdefs.ErrInternal
		}

		token, ok := md["authorization"]
		if len(token) == 0 || !ok {
			log.Error("missing token for ", slog.String("interceptor", op))
			return nil, errdefs.ErrInternal
		}

		res, err := jwtsvc.ValidateToken(ctx, &pbjwtsvc.ValidateTokenRequest{
			Token: token[0],
		})
		if err != nil {
			log.Error("invalid token for ", slog.String("interceptor", op))
			return nil, err
		}

		ctx = context.WithValue(ctx, "uid", res.Uid)

		log.Info("authenticated user: ", res.Uid)

		return handler(ctx, req)
	}
}
