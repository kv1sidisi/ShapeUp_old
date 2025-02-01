package confacchdlr

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
)

// ConfirmAccount interface represents service for confirm account endpoint.
type ConfirmAccount interface {
	ConfirmAccount(token string) error
}

// New creates endpoint for confirm account service.
func New(log *slog.Logger, confirmAccount ConfirmAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.confirm_account"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		tokenString := r.URL.Query().Get("token")

		if tokenString == "" {
			http.Error(w, "Token is missing", http.StatusBadRequest)
			return
		}

		log.Info("got token: ", tokenString)

		if err := confirmAccount.ConfirmAccount(tokenString); err != nil {
			log.Error("Failed confirming account: ", err)
		}
	}
}
