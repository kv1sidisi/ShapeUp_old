package confacchdlr

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
)

// ConfAccSvc interface represents service for confirm account endpoint.
type ConfAccSvc interface {
	ConfirmAccount(token string) error
}

// New creates endpoint for confirm account service.
func New(log *slog.Logger, confirmAccount ConfAccSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.confacchdlr.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tokenString := r.URL.Query().Get("token")

		if tokenString == "" {
			log.Error("token parameter is missing")
			http.Error(w, "token is missing", http.StatusBadRequest)
			return
		}

		log.Info("got token: ", tokenString)

		if err := confirmAccount.ConfirmAccount(tokenString); err != nil {
			log.Error("failed confirming account: ", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
