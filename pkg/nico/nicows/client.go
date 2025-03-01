package nicows

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/websocket-client-go/websocket"
)

type NicoWSClient struct {
	*websocket.WebSocketClient
	ch chan Message
}

func InitPayloads(isReconnecting bool) []Request {
	return []Request{
		{
			Type: START_WATCHING,
			Data: StartWatchingData{
				Reconnect: isReconnecting,
				Room: StartWatching_Room{
					Protocol:    "webSocket",
					Commentable: true,
				},
				Stream: StartWatching_Stream{
					ChasePlay:         false,
					Quality:           "abr",
					Protocol:          "hls+fmp4",
					Latency:           "low",
					AccessRightMethod: "single_cookie",
				},
			},
		},
		{
			Type: GET_AKASHIC,
			Data: GetAkashicData{
				ChasePlay: false,
			},
		},
	}
}

func NewClient(cancel context.CancelFunc) (*NicoWSClient, chan Message) {
	ch := make(chan Message)

	client := &NicoWSClient{
		ch: ch,
		WebSocketClient: websocket.NewWebSocketClient(
			// onOpen
			func(ws *websocket.WebSocketClient, isReconnecting bool) {
				for _, p := range InitPayloads(isReconnecting) {
					log.Println("NicoWS: send init payload type is", p.Type)
					if err := ws.SendJSON(p); err != nil {
						log.Println("NicoWS: failed to send handshake message type of", p.Type)
						go ws.Disconnect(false)
					}
				}
			},

			// onClose
			func(ws *websocket.WebSocketClient, isReconnecting bool) {
				if !isReconnecting {
					close(ch)
					cancel()
				}
			},

			// onMessage
			func(ws *websocket.WebSocketClient, payload []byte) {
				m := utils.Unmarshaller[Message](payload)
				if m == nil {
					return
				}
				log.Println("NicoWS: received type", m.Type)

				switch m.Type {
				default:
					log.Println("NicoWS: unsupported type")

				case STREAM, MESSAGE_SERVER:
					ch <- *m

				case SCHEDULE:
					if data := utils.Unmarshaller[ScheduleData](m.Data); data != nil {
						// Endが"現在時刻" or "以前"なら終了処理
						if data.End.Before(time.Now()) || data.End.Equal(time.Now()) {
							log.Println("NicoWS: program has finished")
							go ws.Disconnect(false)
						}
					}

				case PING:
					for _, t := range []MessageType{PONG, KEEP_SEAT} {
						bytes, _ := json.Marshal(Message{Type: t})
						log.Println("NicoWS: send message", string(bytes))
						if err := ws.SendJSON(Message{Type: t}); err != nil {
							log.Println("NicoWS: failed to send message of", t)
						}
						time.Sleep(500 * time.Millisecond)
					}

				case RECONNECT:
					if data := utils.Unmarshaller[ReconnectData](m.Data); data != nil {
						time.Sleep(time.Duration(data.WaitTimeSec) * time.Second)

						q := ws.URL.Query()
						q.Set("audience_token", data.AudienceToken)
						ws.URL.RawQuery = q.Encode()

						if err := ws.Reconnect(ws.URL.String(), 3, 2); err != nil {
							log.Println("NicoWS: reconnect failed", err)
						}
					}
				}
			},
		),
	}

	return client, ch
}
