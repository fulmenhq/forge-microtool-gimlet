package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fulmenhq/gofulmen/appidentity"
	"github.com/fulmenhq/gofulmen/logging"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads the layered configuration for the application.
// For microtools, config is optional - we provide simple defaults and env var support.
func LoadConfig(ctx context.Context, identity *appidentity.Identity, logger *logging.Logger) (*ExampleConfig, error) {
	// Default config
	cfg := &ExampleConfig{
		InputPath:  ".",
		OutputPath: "./output",
		MaxDepth:   10,
	}

	// Try to load user config file if it exists
	configPath, err := userConfigPath(identity)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(configPath); err == nil {
		data, err := readUserConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		logger.Info("loaded config from file", zap.String("path", configPath))
	}

	// Env var overrides
	if envInput := os.Getenv(identity.EnvPrefix + "INPUT_PATH"); envInput != "" {
		cfg.InputPath = envInput
	}
	if envOutput := os.Getenv(identity.EnvPrefix + "OUTPUT_PATH"); envOutput != "" {
		cfg.OutputPath = envOutput
	}

	logger.Info("configuration loaded successfully")
	return cfg, nil
}

// userConfigPath returns the config file path under the user's .config directory.
func userConfigPath(identity *appidentity.Identity) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	path := filepath.Join(homeDir, ".config", identity.Vendor, identity.ConfigName+".yaml")
	return cleanUserConfigPath(homeDir, path)
}

func cleanUserConfigPath(homeDir, path string) (string, error) {
	root := filepath.Join(homeDir, ".config")
	clean := filepath.Clean(path)
	rel, err := filepath.Rel(root, clean)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("config path is outside the user config directory")
	}
	return clean, nil
}

func readUserConfig(path string) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	clean, err := cleanUserConfigPath(homeDir, path)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(filepath.Join(homeDir, ".config"))
	if err != nil {
		return nil, err
	}
	defer func() { _ = root.Close() }()
	rel, err := filepath.Rel(filepath.Join(homeDir, ".config"), clean)
	if err != nil {
		return nil, err
	}
	f, err := root.Open(rel)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return io.ReadAll(f)
}
