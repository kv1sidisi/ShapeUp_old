package regusrhdlr

import (
	pbusrcreatesvc "GatewayAPI/api/grpccl/pb/usrcreatesvc"
	"encoding/json"
	"github.com/go-chi/chi/middleware"
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
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		resp, err := regUsrSvc.RegisterUser(req.Email, req.Password)
		if err != nil {
			log.Error("failed register account: ", err)
			http.Error(w, "failed to register account", http.StatusInternalServerError)
		}
		log.Info("user registered successfully")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed encoding response: ", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
