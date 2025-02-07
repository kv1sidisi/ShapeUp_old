package regusrhdlr

import (
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	pbusrcreatesvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/usrcreatesvc"
	mapper "github.com/kv1sidisi/shapeup/services/gtwapi/internal/utils/grpchttperrmap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

// RegUsrSvc interface represents service for register user endpoint.
type RegUsrSvc interface {
	RegisterUser(email, password string) (resp *pbusrcreatesvc.RegisterResponse, err error)
}

// JSONRegisterRequest struct for json request parsing.
type JSONRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// New creates endpoint for register user service.
func New(log *slog.Logger, regUsrSvc RegUsrSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "regusrhdlr.register_user"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req JSONRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body: ", err)
			mapper.WriteError(w, status.Error(codes.InvalidArgument, "invalid reuest body"), log)
			return
		}

		resp, err := regUsrSvc.RegisterUser(req.Email, req.Password)
		if err != nil {
			log.Error("failed register account: ", err)
			mapper.WriteError(w, err, log)
			return
		}
		log.Info("user registered successfully")

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", slog.Any("error", err))
			mapper.WriteError(w, status.Error(codes.Internal, "failed to encode response"), log)
			return
		}
	}
}
