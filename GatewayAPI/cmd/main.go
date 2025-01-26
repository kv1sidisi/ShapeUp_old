package main

import (
	"GatewayAPI/internal/config"
	"GatewayAPI/internal/http-server/handlers"
	"GatewayAPI/internal/http-server/middleware/logger"
	"GatewayAPI/internal/service"
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

	//Connection to registration microservice with gRPC
	conn, err := grpc.NewClient(cfg.GRPC.ConfirmAccountAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create grpc client connection to confirm account service: ", err)
		panic(err)
	}
	defer conn.Close()

	log.Info("grpc client connected to", cfg.GRPC.ConfirmAccountAddress)

	client := pb.NewUserCreationClient(conn)

	confirmAccountService := service.New(log, client)

	router.Get("/confirm_account", handlers.New(log, confirmAccountService))

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
