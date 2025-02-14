package recorder

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/shinosaki/namagent/internal/recorder/types"
	"github.com/shinosaki/namagent/utils"
)

var (
	ffmpeg *exec.Cmd
)

func WebSocket(
	sc *utils.SignalContext,
	program *types.ProgramData,
	ffmpegPath string,
	outputPath string,
) error {
	var (
		done         = make(chan struct{})
		webSocketUrl = program.Site.Relive.WebSocketUrl
	)

	ws, _, err := websocket.DefaultDialer.Dial(webSocketUrl, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Send initial payloads
	payloads := []types.WS_Request{
		{
			Type: types.WSRequestType_START_WATCHING,
			Data: types.WS_StartWatching{
				Reconnect: false,
				Room: types.WS_StartWatchingRoom{
					Protocol:    "webSocket",
					Commentable: true,
				},
				Stream: types.WS_StartWatchingStream{
					ChasePlay:         false,
					Quality:           "abr",
					Protocol:          "hls+fmp4",
					Latency:           "low",
					AccessRightMethod: "single_cookie",
				},
			},
		},
		{
			Type: types.WSRequestType_GET_AKASHIC,
			Data: types.WS_GetAkashic{
				ChasePlay: false,
			},
		},
	}
	for _, p := range payloads {
		log.Println("websocket send message:", p.Type)
		ws.WriteJSON(p)
	}

	go onMessage(ws, done, ffmpegPath, outputPath)

	for {
		select {
		case <-done:
			return nil
		case <-sc.Context().Done():
			// Close ffmpeg
			if ffmpeg != nil {
				if err := ffmpeg.Process.Signal(os.Interrupt); err != nil {
					log.Println(err)
				}
			}

			// Close websocket
			ws.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure, "",
				),
			)

			<-done
			return nil
		}
	}
}

func onMessage(
	ws *websocket.Conn,
	done chan struct{},
	ffmpegPath string,
	outputPath string,
) {
	defer close(done)
	var res types.WS_Response

	// Receive message forever
	for {
		if err := ws.ReadJSON(&res); err != nil {
			return
		}

		switch res.Type {
		case types.WSResponseType_PING:
			ws.WriteJSON(types.WS_Request{Type: types.WSRequestType_PONG})
			ws.WriteJSON(types.WS_Request{Type: types.WSRequestType_KEEP_SEAT})
		case types.WSResponseType_STREAM:
			var data types.WS_StreamData
			if err := json.Unmarshal(res.Data, &data); err != nil {
				return
			}

			ffmpeg = FFmpeg(data.URI, data.Cookies, ffmpegPath, outputPath)
			ffmpeg.Run()
		}
	}
}
