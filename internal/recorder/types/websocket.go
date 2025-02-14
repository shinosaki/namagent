package types

import "encoding/json"

type WS_ResponseType string

const (
	WSResponseType_PING   WS_ResponseType = "ping"
	WSResponseType_STREAM WS_ResponseType = "stream"
)

type WS_StreamCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type WS_StreamData struct {
	URI     string            `json:"uri"`
	Cookies []WS_StreamCookie `json:"cookies"`
}

type WS_Stream struct {
	Type WS_ResponseType `json:"type"`
	Data WS_StreamData   `json:"data"`
}

type WS_Ping struct {
	Type WS_ResponseType `json:"type"`
}

type WS_Response struct {
	Type WS_ResponseType `json:"type"`
	Data json.RawMessage `json:"data"`
}

// requests
type WS_RequestType string

const (
	WSRequestType_START_WATCHING WS_RequestType = "startWatching"
	WSRequestType_GET_AKASHIC    WS_RequestType = "getAkashic"
	WSRequestType_PONG           WS_RequestType = "pong"
	WSRequestType_KEEP_SEAT      WS_RequestType = "keepSeat"
)

type WS_Request struct {
	Type WS_RequestType `json:"type"`
	Data interface{}    `json:"data"`
}

type WS_GetAkashic struct {
	ChasePlay bool `json:"chasePlay"`
}

type WS_StartWatchingRoom struct {
	Protocol    string `json:"protocol"`
	Commentable bool   `json:"commentable"`
}

type WS_StartWatchingStream struct {
	ChasePlay         bool   `json:"chasePlay"`
	Quality           string `json:"quality"`
	Protocol          string `json:"protocol"`
	Latency           string `json:"latency"`
	AccessRightMethod string `json:"accessRightMethod"`
}

type WS_StartWatching struct {
	Reconnect bool                   `json:"reconnect"`
	Room      WS_StartWatchingRoom   `json:"room"`
	Stream    WS_StartWatchingStream `json:"stream"`
}
