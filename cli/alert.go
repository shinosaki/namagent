package cli

import (
	"errors"
	"log"
	"sync"

	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/utils/signalcontext"
	"github.com/shinosaki/namagent/pkg/alert/nico"
	"github.com/shinosaki/namagent/pkg/namagent"
	"github.com/spf13/cobra"
)

var alerts = map[string]namagent.Alert{
	"nico": nico.Alert,
}

func Alert() *cobra.Command {
	var wg sync.WaitGroup

	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Live streaming recorder cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, ok := cmd.Context().Value(signalcontext.SELF).(*signalcontext.SignalContext)
			if !ok {
				return errors.New("failed to load SignalContext")
			}
			ctx := sc.Context()

			config, ok := sc.GetValue(signalcontext.CONFIG).(*config.Config)
			if !ok {
				return errors.New("failed to load Config from SignalContext")
			}

			for name, alert := range alerts {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := alert(config, ctx, sc.Cancel); err != nil {
						log.Printf("Alert: %s alert aborted %v", name, err)
					} else {
						log.Printf("Alert: %s alert finished", name)
					}
				}()
			}

			<-ctx.Done()
			wg.Wait()

			log.Println("Alert: finished")
			return nil
		},
	}

	return cmd
}
