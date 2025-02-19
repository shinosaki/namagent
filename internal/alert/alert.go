package alert

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/shinosaki/namagent/internal/consts"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/namagent/pkg/nico/nicoapi"
)

func Alert(
	client *http.Client,
	sc *utils.SignalContext,
) {
	if client == nil {
		client = utils.NewHttp2Client()
	}

	config, _ := sc.GetValue(consts.CONFIG).(*utils.Config)

	ticker := time.NewTicker(config.Alert.CheckIntervalSec * time.Second)
	defer ticker.Stop()

	checkPrograms := func(programs []nicoapi.RecentProgram) {
		for _, program := range programs {
			userId := program.ProgramProvider.Id
			userName := program.ProgramProvider.Name

			// If not followed
			if !slices.Contains(config.Following.Nico, userId) {
				continue
			}

			// If recorded
			if sc.IsActiveTask(program.Id) {
				continue
			}

			log.Println("Alert: detected live stream for user id", userId, userName)
			go func() {
				var proc *exec.Cmd

				sc.AddTask(program.Id, func() {
					// CancelTaskされた場合
					// exec.CommandしたRecorderにSIGINTを送信する
					if proc != nil {
						log.Println("Alert: send interrupt to recorder", program.Id)
						if err := proc.Process.Signal(os.Interrupt); err != nil {
							log.Println("Alert: send interrupt to recorder", program.Id, err)
						}
						if err := proc.Wait(); err != nil {
							log.Println("Alert: terminated recorder process", program.Id, err)
						}
					}
				})
				defer sc.CancelTask(program.Id)

				configPath := sc.GetValue(consts.CONFIG_PATH).(string)
				proc = exec.Command(os.Args[0], "recorder", program.Id, "--config", configPath)
				utils.SetSID(proc)

				if err := proc.Run(); err != nil {
					log.Printf("Alert: %s recorder is failed %s", program.Id, err)
				}
			}()
		}
	}

	fetchPrograms := func(isBulkFetch bool) {
		log.Println("Alert: fetch recent programs...")
		programs, err := nicoapi.FetchRecentPrograms(isBulkFetch, client)
		if err != nil {
			log.Println("Alert: recent programs fetch failed", err)
			return
		}

		log.Printf("Alert: check %d programs", len(programs))
		checkPrograms(programs)
	}

	// bulk fetch in first time
	fetchPrograms(true)

	// monitoring forever
	for {
		select {
		case <-sc.Context().Done():
			log.Println("Alert: receive interrupt...")
			sc.CancelAllTasks()
			return
		case <-ticker.C:
			fetchPrograms(false)
		}
	}
}
