package authhdlr

import (
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	pbauthsvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/authsvc"
	mapper "github.com/kv1sidisi/shapeup/services/gtwapi/internal/utils/grpchttperrmap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

// AuthSvc interface represents service for authentication endpoint.
type AuthSvc interface {
	Login(username, password string) (resp *pbauthsvc.LoginResponse, err error)
}

type JSONAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New creates endpoint for login account service.
func New(log *slog.Logger, authSvc AuthSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.authhdlr.New"
		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req JSONAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body: ", err)
			mapper.WriteError(w, status.Error(codes.InvalidArgument, "invalid reuest body"), log)
			return
		}

		resp, err := authSvc.Login(req.Username, req.Password)
		if err != nil {
			log.Error("failed to log user in", slog.Any("error", err))
			mapper.WriteError(w, err, log)
			return
		}
		log.Info("login succeeded")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", slog.Any("error", err))
			mapper.WriteError(w, status.Error(codes.Internal, "failed to encode response"), log)
			return
		}
	}
}
