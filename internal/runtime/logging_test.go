package runtime

import (
	"testing"

	"github.com/fulmenhq/gofulmen/appidentity"
)

func TestSetupLogging(t *testing.T) {
	identity := &appidentity.Identity{
		BinaryName: "gimlet",
		Vendor:     "fulmenhq",
		EnvPrefix:  "GIMLET_",
		ConfigName: "gimlet",
	}

	logger := SetupLogging(identity)
	if logger == nil {
		t.Fatalf("expected logger, got nil")
	}

	// Should allow basic logging without panic
	logger.Info("logging initialized for test")
}
