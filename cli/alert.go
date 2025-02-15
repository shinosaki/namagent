package cli

import (
	"log"
	"time"

	"github.com/shinosaki/namagent/internal/alert"
	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/utils"
	"github.com/spf13/cobra"
)

func Alert() *cobra.Command {
	var (
		client = utils.NewHttp2Client()
		sc     = utils.NewSignalContext()

		// Load from Config
		configPath string
	)

	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Live streaming alert daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer sc.Wait()

			config, err := config.LoadConfig(configPath)
			if err != nil {
				return err
			}

			// Check all active programs in start up
			if programs, err := alert.FetchRecentPrograms(true, client); err != nil {
				return err
			} else {
				alert.Alert(sc, programs, config)
			}

			// Check programs loop
			ticker := time.NewTicker(config.Meta.FetchInterval * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-sc.Context().Done():
					log.Println("stopping now...")
					sc.Wait()
					return nil
				case <-ticker.C:
					if programs, err := alert.FetchRecentPrograms(false, client); err != nil {
						return err
					} else {
						alert.Alert(sc, programs, config)
					}
				}
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "./config.yaml", "Config file path")

	return cmd
}
