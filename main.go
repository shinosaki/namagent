package main

import (
	"context"

	"github.com/shinosaki/namagent/cli"
	"github.com/shinosaki/namagent/internal/consts"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	var configPath string

	cmd := &cobra.Command{
		Use:   "namagent",
		Short: "Live streaming agent",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			config, _ := utils.LoadConfig(configPath)

			// Create signal context
			sc := utils.NewSignalContext()
			defer sc.Wait()

			// Set value to context
			sc = sc.WithValue(consts.CONFIG, config)
			ctx := context.WithValue(cmd.Context(), consts.SIGNAL_CONTEXT, sc)

			// Set sc to cobra's context
			cmd.SetContext(ctx)

			return nil
		},
	}
	defer cmd.Execute()

	// Sub-commands
	cmd.AddCommand(cli.Alert())
	cmd.AddCommand(cli.Recorder())

	// Global flags
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path of config file, (default: ./config.yaml)")
}
