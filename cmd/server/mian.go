package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mohammadne/snapp-food/cmd"
	"github.com/mohammadne/snapp-food/inernal/api/http"
	"github.com/mohammadne/snapp-food/inernal/config"
	"github.com/mohammadne/snapp-food/inernal/entities"
	"github.com/mohammadne/snapp-food/inernal/repository"
	"github.com/mohammadne/snapp-food/pkg/logger"
)

func main() {
	port := flag.Int("port", 8088, "The server port which handles requests (default: 8088)")
	environmentRaw := flag.String("environment", "", "The environment (default: local)")
	flag.Parse() // Parse the command-line flags

	var cfg config.Config
	var err error

	switch config.ToEnvironment(*environmentRaw) {
	case config.EnvironmentLocal:
		cfg, err = config.LoadDefaults(true)
	default:
		cfg, err = config.Load(true)
	}

	if err != nil {
		log.Fatalf("failed to load config: \n%v", err)
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: \n%v", err)
	}

	logger.Warn("Build Information", cmd.BuildInfo()...)

	bank := repository.NewBank([]entities.User{
		{
			ID:      1,
			Balance: 100,
			Sheba:   "IR7740802513265426484548",
		},
		{
			ID:      2,
			Balance: 100,
			Sheba:   "IR9470104877394934515563",
		},
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go http.New(logger, bank).Serve(ctx, &wg, *port)

	<-ctx.Done()
	wg.Wait()
	logger.Warn("interruption signal recieved, gracefully shutdown the server")
}
