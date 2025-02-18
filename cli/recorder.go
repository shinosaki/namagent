package cli

import (
	"fmt"

	"github.com/shinosaki/namagent/internal/consts"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/namagent/pkg/nico"
	"github.com/spf13/cobra"
)

func Recorder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recorder",
		Short: "Live streaming recorder cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "" {
				return fmt.Errorf("URL or Live ID is required")
			}

			programId := nico.ExtractProgramId(args[0])
			if programId == "" {
				return fmt.Errorf("does not contain programId of argument")
			}

			ctx := cmd.Context()
			sc, ok := ctx.Value(consts.SIGNAL_CONTEXT).(*utils.SignalContext)
			if !ok {
				return fmt.Errorf("failed to load SignalContext")
			}

			return nico.Client(programId, nil, sc)
		},
	}

	return cmd
}
