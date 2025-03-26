package config_test

import (
	"testing"

	"github.com/mohammadne/snapp-food/inernal/config"
)

func TestLoadDefaults(t *testing.T) {
	_, err := config.LoadDefaults(true)
	if err != nil {
		t.Error(err)
	}
}
