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

	"github.com/mohammadne/snapp-food/inernal/api/http/handlers"
	"github.com/mohammadne/snapp-food/inernal/api/http/i18n"
	"github.com/mohammadne/snapp-food/inernal/api/http/middlewares"
	"github.com/mohammadne/snapp-food/inernal/repository"
)

type Server struct {
	logger *zap.Logger
	app    *fiber.App
}

func New(log *zap.Logger, bank repository.Bank) *Server {
	server := &Server{logger: log}

	i18n, err := i18n.New(log)
	if err != nil {
		log.Fatal("failed to load i18n", zap.Error(err))
	}

	server.app = fiber.New(fiber.Config{})

	server.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	handlers.NewHealthz(server.app, log)

	api := server.app.Group("api")
	middlewares.NewLanguage(api, log)
	handlers.NewSheba(api, log, i18n, bank)

	return server
}

func (s *Server) Serve(ctx context.Context, wg *sync.WaitGroup, port int) {
	go func() {
		address := fmt.Sprintf("0.0.0.0:%d", port)

		s.logger.Info("starting server", zap.String("address", address))
		err := s.app.Listen(address)
		s.logger.Fatal("error resolving server", zap.String("address", address), zap.Error(err))
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(shutdownCtx); err != nil {
		s.logger.Error("error shutdown http server", zap.Error(err))
	}

	s.logger.Warn("gracefully shutdown the https servers")
}
