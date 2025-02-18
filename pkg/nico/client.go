package nico

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/shinosaki/namagent/internal/consts"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/namagent/pkg/nico/ndgr"
	"github.com/shinosaki/namagent/pkg/nico/nicowebsocket"
	"github.com/shinosaki/nicolive-comment-protobuf/proto/dwango/nicolive/chat/data"
	"github.com/shinosaki/nicolive-comment-protobuf/proto/dwango/nicolive/chat/service/edge"
)

func Client(
	programId string,
	client *http.Client,
	sc *utils.SignalContext,
) error {
	if client == nil {
		client = utils.NewHttp2Client()
	}

	data, err := FetchProgramData(programId, client)
	if err != nil {
		return err
	}

	if data.Program.Status == ProgramStatus_ENDED {
		return fmt.Errorf("program is ended")
	}

	ws, messageChan := nicowebsocket.NewClient()
	if err := ws.Connect(data.Site.Relive.WebSocketUrl); err != nil {
		return err
	}

	go websocketHandler(sc, ws, client, messageChan, data)

	ws.Wait()
	sc.Wait()

	return nil
}

func websocketHandler(
	sc *utils.SignalContext,
	ws *nicowebsocket.Client,
	client *http.Client,
	messageChan chan nicowebsocket.Message,
	programData *ProgramData,
) {
	sc.AddTask("websocket", func() {})
	defer sc.CancelTask("websocket")

	config, _ := sc.GetValue(consts.CONFIG).(*utils.Config)

	outputBaseName := utils.Template(
		config.Recorder.OutputTemplate,
		map[string]string{
			"yyyymmdd":   time.Now().Format("20060102"),
			"id":         programData.Program.NicoliveProgramId,
			"providerId": programData.Program.Supplier.ProgramProviderId,
			"title":      utils.Escape(programData.Program.Title, ""),
		},
	)
	outputBasePath, err := filepath.Abs(outputBaseName)
	if err != nil {
		log.Panicln("failed to parse output_template", err)
	}

	outputMediaPath := outputBasePath + ".ts"
	outputCommentPath := outputBasePath + ".json"

	for {
		select {
		case <-sc.Context().Done():
			log.Println("NicoClient: receive interrupt...")
			ws.Disconnect()
			return
			// ws.doneの処理が必要かも
		case message, ok := <-messageChan:
			if !ok {
				return
			}

			switch message.Type {
			case nicowebsocket.STREAM:
				var data nicowebsocket.StreamData
				if err := json.Unmarshal(message.Data, &data); err != nil {
					log.Println("NicoClient: unmarshal stream data failed:", err)
					continue
				}
				command := utils.BulkTemplate(
					config.Recorder.CommandTemplate,
					map[string]string{
						"cookies": nicowebsocket.CookieBuilder(data.Cookies),
						"url":     data.URI,
						"output":  outputMediaPath,
					},
				)
				go utils.ExecCommand(command, sc)

			case nicowebsocket.MESSAGE_SERVER:
				var data nicowebsocket.MessageServerData
				if err := json.Unmarshal(message.Data, &data); err != nil {
					log.Println("NicoClient: unmarshal message server data failed:", err)
					continue
				}
				go commentHandler(data.ViewURI, "now", outputCommentPath, sc, client)
			}
		}
	}
}

func commentHandler(
	url string,
	at string,
	outputPath string,
	sc *utils.SignalContext,
	client *http.Client,
) {
	sc.AddTask(url+at, func() {})
	defer sc.CancelTask(url + at)

	var (
		mu              = &sync.Mutex{}
		chatBuffer      []*data.Chat
		alreadySegments = make(map[string]struct{})
	)

	// Periodic write chat data to file
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-sc.Context().Done():
				mu.Lock()
				if err := utils.SaveToFile(chatBuffer, outputPath); err != nil {
					log.Println(err)
				}
				mu.Unlock()
				return
			case <-ticker.C:
				mu.Lock()
				if err := utils.SaveToFile(chatBuffer, outputPath); err != nil {
					log.Println(err)
				}
				chatBuffer = nil
				mu.Unlock()
			}
		}
	}()

	segmentHandler := func(url string) {
		ch, err := ndgr.Reader(url, client, sc)
		if err != nil {
			log.Println("SegmentHandler Error:", err)
			return
		}

		for {
			select {
			case <-sc.Context().Done():
				log.Println("SegmentHandler: receive interrupt...")
				return

			case s, ok := <-ch:
				if !ok {
					return
				}

				if segment, ok := s.(*edge.ChunkedMessage); ok {
					// log.Println("Received Segment:", segment)
					if payload, ok := segment.Payload.(*edge.ChunkedMessage_Message); ok {
						if message, ok := payload.Message.Data.(*data.NicoliveMessage_Chat); ok {
							chat := message.Chat
							log.Println("Chat Message:", chat)
							mu.Lock()
							chatBuffer = append(chatBuffer, chat)
							mu.Unlock()
						}
					}
				}
			}
		}
	}

	for {
		select {
		case <-sc.Context().Done():
			log.Println("CommentHandler: receive interrupt")
			return

		default:
			url := fmt.Sprintf("%s?at=%s", url, at)

			ch, err := ndgr.Reader(url, client, sc)
			if err != nil {
				log.Println(err)
				time.Sleep(2 * time.Second) // 必要？？
				continue
			}

			for chunk := range ch {
				log.Println("Receive chunk:", chunk)

				switch e := chunk.(type) {
				case *edge.ChunkedEntry:
					log.Println("Receive ChunkEntry")

					switch entry := e.Entry.(type) {
					case *edge.ChunkedEntry_Next:
						at = strconv.FormatInt(entry.Next.At, 10)
						log.Println("Next at:", at)

					case *edge.ChunkedEntry_Segment:
						url := entry.Segment.Uri
						if _, exists := alreadySegments[url]; !exists {
							alreadySegments[url] = struct{}{}
							go segmentHandler(url)
						}

						// case *edge.ChunkedEntry_Backward:
						// 	url := entry.Backward.Segment.Uri
						// 	if _, exists := alreadySegments[url]; !exists {
						// 		alreadySegments[url] = struct{}{}
						// 		go segmentHandler(url)
						// 	}

						// case *edge.ChunkedEntry_Previous:

						// case *edge.ChunkedEntry_ReadyForNext:

					}

				case *edge.ChunkedMessage:
					log.Println("Receive ChunkMessage")
					switch payload := e.Payload.(type) {
					case *edge.ChunkedMessage_Message:
						log.Println("ChunkedMessage: message is", payload.Message)

					case *edge.ChunkedMessage_State:
						log.Println("ChunkedMessage: state is", payload.State)

					case *edge.ChunkedMessage_Signal_:
						log.Println("ChunkedMessage: signal is", payload.Signal)

					}

				default:
					log.Println("unknown chunk type")
				}
			}
		}
	}
}
