package cli

import (
	"fmt"

	"github.com/shinosaki/namagent/internal/recorder"
	"github.com/shinosaki/namagent/utils"
	"github.com/spf13/cobra"
)

func Recorder() *cobra.Command {
	var (
		client = utils.NewHttp2Client()
		sc     = utils.NewSignalContext()

		// args
		outputPath string
		ffmpegPath string
	)

	cmd := &cobra.Command{
		Use:   "recorder",
		Short: "Live streaming recording cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer sc.Wait()

			if len(args) == 0 || args[0] == "" {
				return fmt.Errorf("url or liveid is require")
			}
			if outputPath == "" {
				return fmt.Errorf("output path is require")
			}

			programId := recorder.ExtractProgramId(args[0])
			if programId == "" {
				return fmt.Errorf("does not contain program id of arg")
			}

			return recorder.Recorder(sc, client, programId, ffmpegPath, outputPath)
		},
	}

	cmd.Flags().StringVar(&ffmpegPath, "ffmpeg", "ffmpeg", "FFmpeg Path (Default: 'ffmpeg')")
	cmd.Flags().StringVar(&outputPath, "output", "", "Output Path")

	return cmd
}
