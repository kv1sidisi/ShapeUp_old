package authhdlr

import (
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	pbauthsvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/authsvc"
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

func New(log *slog.Logger, authSvc AuthSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.authhdlr.New"
		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req JSONAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		resp, err := authSvc.Login(req.Username, req.Password)
		if err != nil {
			log.Error("failed log user in: ", err)
			http.Error(w, "failed to log user in ", http.StatusInternalServerError)
		}
		log.Info("login succeeded")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response: ", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
