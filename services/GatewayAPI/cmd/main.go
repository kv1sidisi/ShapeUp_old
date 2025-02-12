package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	loadconfig "github.com/kv1sidisi/shapeup/pkg/config"
	"github.com/kv1sidisi/shapeup/pkg/logger"
	pbauthsvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/authsvc"
	pbusrcreatesvc "github.com/kv1sidisi/shapeup/services/gtwapi/api/grpccl/pb/usrcreatesvc"
	"github.com/kv1sidisi/shapeup/services/gtwapi/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/gtwapi/cmd/grpccl/consts"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/config"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/http-server/handlers/authhdlr"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/http-server/handlers/confacchdlr"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/http-server/handlers/regusrhdlr"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/http-server/middleware/midlogger"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/service/authsvc"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/service/confaccsvc"
	"github.com/kv1sidisi/shapeup/services/gtwapi/internal/service/regusrsvc"
	"log/slog"
	"net/http"
)

func main() {
	cfg := &config.Config{}
	loadconfig.MustLoad(cfg)

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(midlogger.New(log))
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
