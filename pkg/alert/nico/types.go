package nico

import "encoding/json"

type Meta struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"errorCode"`
}

type APIResponse struct {
	Meta Meta            `json:"meta"`
	Data json.RawMessage `json:"data"`
}

type (
	ProviderType string
	LiveCycle    string
)

const (
	ProviderType_COMMUNITY ProviderType = "community"
	LiveCycle_ONAIR        LiveCycle    = "ON_AIR"
	LiveCycle_ENDED        LiveCycle    = "ENDED"
)

type RecentProgram struct {
	Id    string `json:"id"`
	Title string `json:"title"`

	ListingThumbnail        string `json:"listingThumbnail"`
	FlippedListingThumbnail string `json:"flippedListingThumbnail"`
	WatchPageUrl            string `json:"watchPageUrl"`
	WatchPageUrlAtExtPlayer string `json:"watchPageUrlAtExtPlayer"`

	ProviderType ProviderType `json:"providerType"`
	LiveCycle    LiveCycle    `json:"liveCycle"`

	// 13-digit timestamp
	BeginAt int `json:"beginAt"`
	EndAt   int `json:"endAt"`

	IsFollowerOnly              bool `json:"isFollowerOnly"`
	IsPayProgram                bool `json:"isPayProgram"`
	IsOfficialChannelMemberFree bool `json:"isOfficialChannelMemberFree"`

	ProgramProvider ProgramProvider `json:"programProvider"`
	SocialGroup     SocialGroup     `json:"socialGroup"`
	Statistics      Statistics      `json:"statistics"`
	Timeshift       Timeshift       `json:"timeshift"`
}

// 配信者情報
type ProgramProvider struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	IconSmall string `json:"iconSmall"`
}

// 旧コミュニティ
type SocialGroup struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type Statistics struct {
	WatchCount       int `json:"watchCount"`
	CommentCount     int `json:"commentCount"`
	ReservationCount int `json:"reservationCount"`
}

type Timeshift struct {
	IsPlayable   bool `json:"isPlayable"`
	IsReservable bool `json:"isReservable"`
}
