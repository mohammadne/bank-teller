package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/mohammadne/bank-teller/pkg/logger"
)

type Config struct {
	Logger *logger.Config `required:"true"`
}

func Load(environment Environment) (config Config, err error) {
	switch environment {
	case EnvironmentLocal:
		if err := setDefaults(); err != nil {
			return Config{}, fmt.Errorf("error set default config: %v", err)
		}
	}

	prefix := strings.ToUpper(strings.ReplaceAll(System, "-", "_"))

	if err = envconfig.Process(prefix, &config); err != nil {
		return Config{}, fmt.Errorf("error processing config via envconfig: %v", err)
	}

	fmt.Println("================ Loaded Configuration ================")
	object, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(object))
	fmt.Println("======================================================")

	return config, nil
}

const seperator = "_"

//go:embed defaults.env
var defaults string

func setDefaults() error {
	lines := strings.Split(defaults, "\n")
	for _, line := range lines {
		splits := strings.Split(line, "=")
		if len(splits) < 2 {
			continue
		}

		key := strings.ReplaceAll(splits[0], seperator+seperator, seperator)
		err := os.Setenv(key, splits[1])
		if err != nil {
			return fmt.Errorf("error set environment %s: %v", key, err)
		}
	}

	return nil
}
