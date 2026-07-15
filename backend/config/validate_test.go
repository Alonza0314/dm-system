package config_test

import (
	"backend/config"
	"testing"

	"github.com/free-ran-ue/util"
)

func TestConfigValidation(t *testing.T) {
	systemConfig := config.Config{}

	if err := util.LoadFromYaml("../../config.yaml", &systemConfig); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if err := systemConfig.Validate(); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}
}
