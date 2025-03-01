package recorder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/namagent/pkg/namagent"
	"github.com/shinosaki/namagent/pkg/nico"
)

var plugins = map[string]*namagent.Plugin{
	"nico": nico.Plugin,
}

func Recorder(
	url string,
	config *config.Config,
	ctx context.Context,
	cancel context.CancelFunc,
) error {
	var (
		wg         sync.WaitGroup
		programId  string
		siteName   string
		sitePlugin *namagent.Plugin
	)

	// Match url
	for name, plugin := range plugins {
		id := plugin.ExtractId(url)

		if id != "" {
			siteName = name
			programId = id
			sitePlugin = plugin

			log.Println("Plugin:", siteName)
			log.Println("Program ID:", id)
			break
		}
	}
	if programId == "" {
		return errors.New("unmatch any plugin")
	}

	// Run client
	commentChan, streamDataChan, err := sitePlugin.Client(
		programId,
		config,
		sitePlugin.Session(config),
		ctx,
	)
	if err != nil {
		return fmt.Errorf("%s Plugin Error: %v", siteName, err)
	}
	log.Println("Recorder: successful connected Plugin Client")

	streamData := <-streamDataChan
	streamData.Extension = config.Recorder.Extension

	// Output Template
	streamData.Output, err = utils.OutputTemplate(programId, config.Recorder.OutputTemplate, streamData.Template)
	if err != nil {
		log.Println("failed to build output template, used default name:", err)
		streamData.Output = time.Now().Format("20060102") + "-" + programId
	}

	streamData.Output = utils.Escape(streamData.Output, "")

	// Unique Filename
	streamData.Output = strings.TrimSuffix(
		utils.UniqueFilename(streamData.Output+"."+streamData.Extension),
		"."+streamData.Extension,
	)

	log.Println("Recorder: output name is", streamData.Output)

	// Create Output Dir
	if outputDir, err := utils.MkDir(streamData.Output); err != nil {
		return fmt.Errorf("failed to create output directory: %v, %s", err, outputDir)
	}

	// Exec Command
	wg.Add(1)
	go func() {
		defer wg.Done()

		// コマンド生成
		command, err := utils.BulkOutputTemplate(programId+"cmd", config.Recorder.CommandTemplate, streamData)
		if err != nil {
			log.Printf("failed to build command template: %v", err)
			return
		}

		proc := exec.Command(command[0], command[1:]...)

		// コマンド実行
		if err := proc.Start(); err != nil {
			log.Printf("failed to exec %s: %v [%s]", command[0], err, strings.Join(command, " "))
			return
		}

		// 終了待機
		done := make(chan error, 1)
		go func() { done <- proc.Wait() }()

		select {
		case <-ctx.Done():
			log.Println("Recorder: receive interrupt")

			// procにSIGINTを送信
			if proc.Process != nil {
				log.Println("Recorder: send SIGINT to command")
				if err := proc.Process.Signal(os.Interrupt); err != nil {
					log.Println("SIGINT output of command:", err)
				}
			}
			<-done // プロセスの終了を待機

		case err := <-done:
			if err != nil {
				log.Printf("%s aborted: %v", command[0], err)
			} else {
				log.Printf("%s finished: %v", command[0], err)
			}
		}
	}()

	// Save comment
	wg.Add(1)
	go func() {
		defer wg.Done()
		var buffer []any

		save := func() (ok bool) {
			if len(buffer) > 0 {
				if err := saveCommentToJSON(streamData.Output, buffer); err != nil {
					log.Println("failed to write comment:", err)
					return false
				}
			}
			return true
		}

		for {
			select {
			case <-ctx.Done():
				// interruptを受信したら保存して終了
				log.Println("commenthandler: receive interrupt")
				save()
				return

			case comment, ok := <-commentChan:
				// チャンネルがcloseしたら保存して終了
				if !ok {
					save()
					return
				}

				// bufferが10以上で保存
				buffer = append(buffer, comment)
				if len(buffer) >= 10 {
					if ok := save(); !ok {
						// 保存に失敗したらbufferを維持したまま再試行
						continue
					}
					// 保存に成功したらbufferをクリア
					buffer = nil
				}
			}
		}
	}()

	// CommentHandler, ExecCommand Goroutinesの終了を待機
	go func() {
		wg.Wait()
		cancel()
	}()
	<-ctx.Done()
	wg.Wait()

	log.Println("Recorder: finished")
	return nil
}

func saveCommentToJSON(output string, buffer []any) error {
	path := output + ".json"

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	var exists []any
	if err := json.NewDecoder(file).Decode(&exists); err != nil && err != io.EOF {
		return fmt.Errorf("error decoding existing data: %v", err)
	}
	exists = append(exists, buffer...)

	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exists); err != nil {
		return fmt.Errorf("error encoding data to JSON: %v", err)
	}

	log.Printf("Saved %d items to %s", len(buffer), output)
	return nil
}
