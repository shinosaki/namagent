package nicowebsocket

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	// request type
	START_WATCHING MessageType = "startWatching"
	GET_AKASHIC    MessageType = "getAkashic"

	// ping type
	PING      MessageType = "ping"
	PONG      MessageType = "pong"
	KEEP_SEAT MessageType = "keepSeat"

	// message type
	STREAM         MessageType = "stream"
	RECONNECT      MessageType = "reconnect"
	MESSAGE_SERVER MessageType = "messageServer"
)

type Request struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type GetAkashicData struct {
	ChasePlay bool `json:"chasePlay"`
}

type StartWatchingData struct {
	Reconnect bool                 `json:"reconnect"`
	Room      StartWatching_Room   `json:"room"`
	Stream    StartWatching_Stream `json:"stream"`
}

type StartWatching_Room struct {
	Protocol    string `json:"protocol"`
	Commentable bool   `json:"commentable"`
}

type StartWatching_Stream struct {
	ChasePlay         bool   `json:"chasePlay"`
	Quality           string `json:"quality"`
	Protocol          string `json:"protocol"`
	Latency           string `json:"latency"`
	AccessRightMethod string `json:"accessRightMethod"`
}

type ReconnectData struct {
	AudienceToken string `json:"audienceToken"`
	WaitTimeSec   int    `json:"WaitTimeSec"`
}

type StreamData struct {
	URI                string   `json:"uri"`
	SyncUri            string   `json:"syncUri"`
	Quality            string   `json:"quality"`
	Protocol           string   `json:"protocol"`
	AvailableQualities []string `json:"availableQualities"`
	Cookies            []Cookie `json:"cookies"`
}

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type MessageServerData struct {
	ViewURI      string    `json:"viewUri"`
	VposBaseTime time.Time `json:"vposBaseTime"`
}
