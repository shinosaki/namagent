package alert

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"time"

	"github.com/shinosaki/namagent/internal/alert/types"
	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/recorder"
	"github.com/shinosaki/namagent/utils"
)

func Alert(
	sc *utils.SignalContext,
	programs []types.RecentProgram,
	config config.Config,
) {
	var (
		ffmpegPath     = config.Paths.FFmpeg
		followingUsers = config.Following.Users["nico"]
	)

	for _, program := range programs {
		if slices.Contains(followingUsers, program.ProgramProvider.Id) {
			// show information
			log.Printf("Live stream detected for followed user: %s (%s)",
				program.ProgramProvider.Id,
				program.ProgramProvider.Name,
			)

			// active recording state management
			if sc.IsActiveTask(program.Id) {
				log.Println("This program is already recording:", program.Id)
				continue
			}

			_, cancel := context.WithCancel(sc.Context())
			sc.AddTask(program.Id, cancel)

			go func() {
				defer sc.CancelTask(program.Id)

				outputPath, err := filepath.Abs(
					filepath.Join(
						config.Paths.OutputBaseDir,
						program.ProgramProvider.Id,
						fmt.Sprintf("%s-%s-%s-%s.%s",
							time.Now().Format("20060102"),
							program.Id,
							program.ProgramProvider.Name,
							program.Title,
							"ts",
						),
					),
				)
				if err != nil {
					log.Println(err)
					return
				}

				log.Println("Output Path:", outputPath)

				// run recorder
				recorder.Recorder(sc, nil, program.Id, ffmpegPath, outputPath)
			}()
		}
	}
}
