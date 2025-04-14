package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mohammadne/bank-teller/cmd"
	"github.com/mohammadne/bank-teller/inernal/api/http"
	"github.com/mohammadne/bank-teller/inernal/config"
	"github.com/mohammadne/bank-teller/inernal/entities"
	"github.com/mohammadne/bank-teller/inernal/repository"
	"github.com/mohammadne/bank-teller/inernal/usecases"
	"github.com/mohammadne/bank-teller/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	monitorPort := flag.Int("monitor-port", 8001, "The server port which handles monitoring endpoints (default: 8001)")
	requestPort := flag.Int("request-port", 8002, "The server port which handles http requests (default: 8002)")
	environmentRaw := flag.String("environment", "", "The environment (default: local)")
	flag.Parse() // Parse the command-line flags

	environment := config.ToEnvironment(*environmentRaw)
	cfg, err := config.Load(environment)
	if err != nil {
		log.Panicf("failed to load config: \n%v", err)
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: \n%v", err)
	}

	buildInformations := make([]zap.Field, 0)
	for key, value := range cmd.BuildInfo() {
		buildInformations = append(buildInformations, zap.String(key, value))
	}
	logger.Warn("Build Information", buildInformations...)

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

	// usecases
	sheba := usecases.NewSheba(bank)
	users := usecases.NewUsers(bank)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go http.New(logger, sheba, users).Serve(ctx, &wg, *monitorPort, *requestPort)

	<-ctx.Done()
	wg.Wait()
	logger.Warn("interruption signal recieved, gracefully shutdown the server")
}
