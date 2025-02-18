package cli

import (
	"fmt"

	"github.com/shinosaki/namagent/internal/alert"
	"github.com/shinosaki/namagent/internal/consts"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/spf13/cobra"
)

func Alert() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Live streaming alert daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			sc, ok := ctx.Value(consts.SIGNAL_CONTEXT).(*utils.SignalContext)
			if !ok {
				return fmt.Errorf("failed to load SignalContext")
			}

			alert.Alert(nil, sc)
			return nil
		},
	}

	return cmd
}
