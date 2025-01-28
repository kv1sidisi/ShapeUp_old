package main

import (
	"GatewayAPI/internal/config"
	"GatewayAPI/internal/http-server/handlers/confirm_account"
	handler_register_user "GatewayAPI/internal/http-server/handlers/register_user"
	"GatewayAPI/internal/http-server/middleware/logger"
	"GatewayAPI/internal/service/confirm_account"
	service_register_user "GatewayAPI/internal/service/register_user"
	pb "GatewayAPI/pkg/grpc_client/pb"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	log.Info("connecting to grpc RegistrationService", slog.String("address", cfg.GRPC.ConfirmAccountAddress))
	conn, err := grpc.NewClient(cfg.GRPC.ConfirmAccountAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create grpc client connection to confirm account service: ", err)
		panic(err)
	}
	defer conn.Close()
	log.Info("grpc client connected")
	client := pb.NewUserCreationClient(conn)

	confirmAccountService := service_confirm_account.New(log, client)
	log.Info("confirm account service registered")
	router.Get("/confirm_account", handler_confirm_account.New(log, confirmAccountService))
	log.Info("confirm_account endpoint registered")

	registerUserService := service_register_user.New(log, client)
	log.Info("register user service registered")
	router.Get("/register_user", handler_register_user.New(log, registerUserService))
	log.Info("register user endpoint registered")

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
