package main

import (
	"github.com/shinosaki/namagent/cli"
	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/utils/signalcontext"
	"github.com/spf13/cobra"
)

func main() {
	var configPath string

	cmd := &cobra.Command{
		Use:   "namagent",
		Short: "Live streaming agent",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config, err := config.LoadConfig(configPath)
			if err != nil {
				return err
			}

			sc := signalcontext.NewSignalContext()
			sc = sc.WithValue(signalcontext.CONFIG, config)

			cmd.SetContext(sc.SelfContext())
			return nil
		},
	}
	defer cmd.Execute()

	// Sub-commands
	cmd.AddCommand(cli.Alert())
	cmd.AddCommand(cli.Recorder())

	// Global flags
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path of config file. by default './config.yaml'")
}
