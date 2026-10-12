package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/fulmenhq/gofulmen/appidentity"
	"github.com/fulmenhq/gofulmen/logging"
)

func TestLoadConfig_DefaultsAndEnvOverrides(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	identity := &appidentity.Identity{
		BinaryName: "gimlet",
		Vendor:     "acme",
		EnvPrefix:  "TEST_",
		ConfigName: "tool",
	}
	logger, err := logging.NewCLI(identity.BinaryName)
	if err != nil {
		t.Fatalf("logger init failed: %v", err)
	}

	// Env should override defaults
	t.Setenv("TEST_INPUT_PATH", "/env/input")
	t.Setenv("TEST_OUTPUT_PATH", "/env/output")

	cfg, err := LoadConfig(context.Background(), identity, logger)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.InputPath != "/env/input" {
		t.Fatalf("expected env override for InputPath, got %s", cfg.InputPath)
	}
	if cfg.OutputPath != "/env/output" {
		t.Fatalf("expected env override for OutputPath, got %s", cfg.OutputPath)
	}
	if cfg.MaxDepth != 10 {
		t.Fatalf("expected default MaxDepth 10, got %d", cfg.MaxDepth)
	}
}

func TestLoadConfig_FileAndEnvPriority(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	identity := &appidentity.Identity{
		BinaryName: "gimlet",
		Vendor:     "vendorx",
		EnvPrefix:  "MYTOOL_",
		ConfigName: "mytool",
	}
	logger, err := logging.NewCLI(identity.BinaryName)
	if err != nil {
		t.Fatalf("logger init failed: %v", err)
	}

	// Prepare config file
	configPath, err := userConfigPath(identity)
	if err != nil {
		t.Fatalf("userConfigPath returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	fileData := []byte("input_path: /file/input\noutput_path: /file/output\nmax_depth: 3\n")
	if err := os.WriteFile(configPath, fileData, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Env overrides only inputPath; others should come from file
	t.Setenv("MYTOOL_INPUT_PATH", "/env/input")

	cfg, err := LoadConfig(context.Background(), identity, logger)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.InputPath != "/env/input" {
		t.Fatalf("expected env override for InputPath, got %s", cfg.InputPath)
	}
	if cfg.OutputPath != "/file/output" {
		t.Fatalf("expected OutputPath from file, got %s", cfg.OutputPath)
	}
	if cfg.MaxDepth != 3 {
		t.Fatalf("expected MaxDepth from file (3), got %d", cfg.MaxDepth)
	}
}

func TestUserConfigPath_RejectsEscape(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	outside := filepath.Join(home, "secret.yaml")
	if err := os.WriteFile(outside, []byte("input_path: /nope\n"), 0o644); err != nil {
		t.Fatalf("write outside file: %v", err)
	}

	identity := &appidentity.Identity{
		BinaryName: "gimlet",
		Vendor:     "..",
		EnvPrefix:  "TEST_",
		ConfigName: "secret",
	}
	if _, err := userConfigPath(identity); err == nil {
		t.Fatal("expected config path outside ~/.config to be rejected")
	}
	if _, err := readUserConfig(outside); err == nil {
		t.Fatal("expected a read outside ~/.config to be rejected")
	}
}

func TestReadUserConfig_RejectsSymlinkEscape(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configDir := filepath.Join(home, ".config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	outside := filepath.Join(home, "secret.yaml")
	if err := os.WriteFile(outside, []byte("input_path: /nope\n"), 0o644); err != nil {
		t.Fatalf("write outside file: %v", err)
	}
	link := filepath.Join(configDir, "escape.yaml")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	if _, err := readUserConfig(link); err == nil {
		t.Fatal("symlink under ~/.config pointing outside was accepted")
	}
}
