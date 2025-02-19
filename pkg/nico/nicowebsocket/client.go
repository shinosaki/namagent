package nicowebsocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/shinosaki/namagent/internal/websocket"
)

type Client struct {
	*websocket.Client
	messageChan chan Message
}

func NewClient() (*Client, chan Message) {
	messageChan := make(chan Message)

	initPayloads := []Request{
		{
			Type: START_WATCHING,
			Data: StartWatchingData{
				Reconnect: false,
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

	client := &Client{
		messageChan: messageChan,
		Client: websocket.NewClient(
			func(c *websocket.Client) {
				// onOpen
				for _, payload := range initPayloads {
					if err := c.SendJSON(payload); err != nil {
						log.Println("Nico Websocket: failed to send init message", payload.Type)
						c.Disconnect()
					}
				}
			},
			func(c *websocket.Client) {
				// onClose
				close(messageChan)
			},
			func(c *websocket.Client, payload []byte) {
				// onMessage
				var message Message
				if err := json.Unmarshal(payload, &message); err != nil {
					log.Println("Nico WebSocket: failed to unmarshal payload", err)
					return
				}

				log.Println("Nico WebSocket: received type is", message.Type)
				// log.Println("Nico WebSocket: received data", string(message.Data))

				switch message.Type {
				default:
					log.Println("Nico WebSocket: unsupported type")

				case STREAM, MESSAGE_SERVER:
					messageChan <- message

				case PING:
					for _, t := range []MessageType{PONG, KEEP_SEAT} {
						if err := c.SendJSON(Message{Type: t}); err != nil {
							log.Println("Nico WebSocket: failed to send", t)
						}
					}

				case SCHEDULE:
					var data ScheduleData
					if err := json.Unmarshal(message.Data, &data); err != nil {
						log.Println("Nico Websocket: unmarshal schedule data failed", err)
						return
					}

					// End が現在時刻もしくは以前なら終了処理
					if data.End.Before(time.Now()) || data.End.Equal(time.Now()) {
						log.Println("Nico Websocket: received program has finished")
						c.Disconnect()
					}

				case RECONNECT:
					var data ReconnectData
					if err := json.Unmarshal(message.Data, &data); err != nil {
						log.Println("Nico WebSocket: unmarshal reconnect data failed", err)
						return
					}

					time.Sleep(time.Duration(data.WaitTimeSec) * time.Second)

					// replace audience token of websocket url
					q := c.URL.Query()
					q.Set("audience_token", data.AudienceToken)
					c.URL.RawQuery = q.Encode()

					if err := c.Reconnect(c.URL.String()); err != nil {
						log.Println("Nico WebSocket: reconnect failed", err)
					}
				}
			},
		),
	}

	return client, messageChan
}
