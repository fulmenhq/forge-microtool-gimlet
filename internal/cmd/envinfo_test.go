package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/fulmenhq/gofulmen/appidentity"
	"github.com/fulmenhq/gofulmen/logging"

	"github.com/fulmenhq/forge-microtool-gimlet/internal/config"
)

func TestEnvinfoCommand_Executes(t *testing.T) {
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

	cmd := newEnvinfoCmd(identity)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("envinfo command returned error: %v", err)
	}
}
