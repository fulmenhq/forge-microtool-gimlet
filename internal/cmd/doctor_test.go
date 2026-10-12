package cmd

import (
	"context"
	"testing"

	"github.com/fulmenhq/gofulmen/appidentity"
	"github.com/fulmenhq/gofulmen/logging"

	"github.com/fulmenhq/forge-microtool-gimlet/internal/config"
)

func TestDoctorCommand_Succeeds(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	identity := &appidentity.Identity{
		BinaryName:  "gimlet",
		Vendor:      "fulmenhq",
		EnvPrefix:   "GIMLET_",
		ConfigName:  "gimlet",
		Description: "Test tool",
	}
	logger, err := logging.NewCLI(identity.BinaryName)
	if err != nil {
		t.Fatalf("logger init failed: %v", err)
	}
	cfg := &config.ExampleConfig{}
	Initialize(context.Background(), identity, logger, cfg)

	cmd := newDoctorCmd(identity)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("doctor command returned error: %v", err)
	}
}
