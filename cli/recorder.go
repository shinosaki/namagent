package cli

import (
	"errors"
	"fmt"

	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/recorder"
	"github.com/shinosaki/namagent/internal/utils/signalcontext"
	"github.com/spf13/cobra"
)

func Recorder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recorder",
		Short: "Live streaming recorder cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, ok := cmd.Context().Value(signalcontext.SELF).(*signalcontext.SignalContext)
			if !ok {
				return errors.New("failed to load SignalContext")
			}

			config, ok := sc.GetValue(signalcontext.CONFIG).(*config.Config)
			if !ok {
				return errors.New("failed to load Config from SignalContext")
			}

			url := args[0]
			if url == "" {
				return fmt.Errorf("URL or ProgramID is required")
			}

			return recorder.Recorder(url, config, sc.Context(), sc.Cancel)
		},
	}

	return cmd
}
