package main

import (
	pbauthsvc "GatewayAPI/api/grpccl/pb/authsvc"
	pbusrcreatesvc "GatewayAPI/api/grpccl/pb/usrcreatesvc"
	"GatewayAPI/cmd/grpccl"
	"GatewayAPI/cmd/grpccl/consts"
	"GatewayAPI/internal/config"
	"GatewayAPI/internal/http-server/handlers/authhdlr"
	"GatewayAPI/internal/http-server/handlers/confacchdlr"
	"GatewayAPI/internal/http-server/handlers/regusrhdlr"
	"GatewayAPI/internal/http-server/middleware/logger"
	"GatewayAPI/internal/service/authsvc"
	"GatewayAPI/internal/service/confaccsvc"
	"GatewayAPI/internal/service/regusrsvc"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//connecting grpc clients
	clients := grpccl.New(log, cfg)
	defer clients.Close()

	confAccSvc := confaccsvc.New(log, clients.Cl[consts.UsrCreateSvc].Client.(pbusrcreatesvc.UserCreationClient))
	log.Info("confirm account service registered")
	router.Get("/confirm_account", confacchdlr.New(log, confAccSvc))
	log.Info("confirm_account endpoint registered")

	regUsrSvc := regusrsvc.New(log, clients.Cl[consts.ConfAccSvc].Client.(pbusrcreatesvc.UserCreationClient))
	log.Info("register user service registered")
	router.Post("/register_user", regusrhdlr.New(log, regUsrSvc))
	log.Info("register user endpoint registered")

	authSvc := authsvc.New(log, clients.Cl[consts.AuthSvc].Client.(pbauthsvc.AuthClient))
	log.Info("authentication service registered")
	router.Get("/login", authhdlr.New(log, authSvc))
	log.Info("authentication endpoint registered")

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("address", cfg.Address))
	}
}

// setupLogger initializes logger dependent on environment
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
