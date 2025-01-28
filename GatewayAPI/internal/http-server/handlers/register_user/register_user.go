package handler_register_user

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
)

type RegisterUser interface {
	RegisterUser(email, password string) error
}

func New(log *slog.Logger, registerUser RegisterUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.register_user"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//TODO: parse request

		log.Info("got email: ", email)
		log.Info("got password: ", password)

		if err := registerUser.RegisterUser(email, password); err != nil {
			log.Error("Failed confirming account: ", err)
		}
	}
}
