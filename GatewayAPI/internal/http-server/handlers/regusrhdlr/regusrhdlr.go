package regusrhdlr

import (
	pbusrcreatesvc "GatewayAPI/api/grpccl/pb/usrcreatesvc"
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
)

// RegisterUser interface represents service for register user endpoint.
type RegisterUser interface {
	RegisterUser(email, password string) (resp *pbusrcreatesvc.RegisterResponse, err error)
}

type JSONRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// New creates endpoint for register user service.
func New(log *slog.Logger, registerUser RegisterUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.register_user"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req JSONRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		log.Info("got email: ", req.Email)
		log.Info("got password: ", req.Password)

		resp, err := registerUser.RegisterUser(req.Email, req.Password)
		if err != nil {
			log.Error("Failed confirming account: ", err)
		}
		log.Info("user registered successfully")
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		//TODO: handle errors. look for better practice
		if err != nil {
			return
		}
	}
}
