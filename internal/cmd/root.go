package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/fulmenhq/gofulmen/appidentity"
	"github.com/fulmenhq/gofulmen/logging"
	"github.com/spf13/cobra"

	"github.com/fulmenhq/forge-microtool-gimlet/internal/config"
)

var (
	// Global flags
	logLevel   string
	configPath string
)

// rootCmd represents the base command
var rootCmd *cobra.Command

// loggerInstance holds the logger (used by subcommands)
var loggerInstance *logging.Logger

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// Initialize sets up the root command with app identity, logger, and config
func Initialize(ctx context.Context, identity *appidentity.Identity, logger *logging.Logger, cfg *config.ExampleConfig) {
	loggerInstance = logger

	rootCmd = &cobra.Command{
		Use:   identity.BinaryName,
		Short: identity.Description,
		Long:  identity.Description,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set log level from flag if provided
			if logLevel != "" {
				// Parse level: trace, debug, info, warn, error
				switch strings.ToUpper(logLevel) {
				case "TRACE":
					logger.SetLevel(logging.TRACE)
				case "DEBUG":
					logger.SetLevel(logging.DEBUG)
				case "INFO":
					logger.SetLevel(logging.INFO)
				case "WARN", "WARNING":
					logger.SetLevel(logging.WARN)
				case "ERROR":
					logger.SetLevel(logging.ERROR)
				default:
					return fmt.Errorf("invalid log level: %s (use trace|debug|info|warn|error)", logLevel)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "logging level (trace|debug|info|warn|error)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path override")

	// Add subcommands
	rootCmd.AddCommand(newVersionCmd(identity))
	rootCmd.AddCommand(newEnvinfoCmd(identity))
	rootCmd.AddCommand(newDoctorCmd(identity))
	rootCmd.AddCommand(newExampleCmd(identity))
	rootCmd.AddCommand(newValidateCmd(identity))
}
