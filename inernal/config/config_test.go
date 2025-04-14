package config_test

import (
	"testing"

	"github.com/mohammadne/bank-teller/inernal/config"
)

func TestLoadDefaults(t *testing.T) {
	_, err := config.Load(config.EnvironmentLocal)
	if err != nil {
		t.Error(err)
	}
}
