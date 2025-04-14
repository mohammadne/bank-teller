package http

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/mohammadne/bank-teller/inernal/api/http/handlers"
	"github.com/mohammadne/bank-teller/inernal/api/http/i18n"
	"github.com/mohammadne/bank-teller/inernal/api/http/middlewares"
	"github.com/mohammadne/bank-teller/inernal/repository"
)

type Server struct {
	logger *zap.Logger

	monitorApp *fiber.App
	requestApp *fiber.App
}

func New(log *zap.Logger, bank repository.Bank) *Server {
	server := &Server{logger: log}

	{
		server.monitorApp = fiber.New(fiber.Config{})

		server.monitorApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
		handlers.NewHealthz(server.monitorApp, log)
	}

	i18n, err := i18n.New(log)
	if err != nil {
		log.Fatal("failed to load i18n", zap.Error(err))
	}

	{
		server.requestApp = fiber.New(fiber.Config{})

		api := server.requestApp.Group("api")
		middlewares.NewLanguage(api, log)
		handlers.NewSheba(api, log, i18n, bank)
	}

	return server
}

func (s *Server) Serve(ctx context.Context, wg *sync.WaitGroup, monitorPort, requestPort int) {
	defer wg.Done()

	servers := map[*fiber.App]int{
		s.monitorApp: monitorPort,
		s.requestApp: requestPort,
	}

	for app, port := range servers {
		go func() {
			address := fmt.Sprintf("0.0.0.0:%d", port)
			s.logger.Info("starting server", zap.String("address", address))
			err := app.Listen(address, fiber.ListenConfig{DisableStartupMessage: true})
			s.logger.Fatal("error resolving server", zap.String("address", address), zap.Error(err))
		}()
	}

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for app := range servers {
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			s.logger.Error("error shutdown http server", zap.Error(err))
		}
	}

	s.logger.Warn("gracefully shutdown the https servers")
}
