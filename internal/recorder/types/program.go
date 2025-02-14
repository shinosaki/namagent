package types

type StreamType string

const (
	StreamType_DMC   StreamType = "dmc"
	StreamType_DLIVE StreamType = "dlive"
)

type Stream struct {
	Type StreamType `json:"type"`
}

type Relive struct {
	ApiBaseUrl        string `json:"apiBaseUrl"`
	ChannelApiBaseUrl string `json:"channelApiBaseUrl"`
	WebSocketUrl      string `json:"webSocketUrl"`
	CsrfToken         string `json:"csrfToken"`
	AudienceToken     string `json:"audienceToken"`
}

type Site struct {
	Relive     Relive `json:"relive"`
	FrontendId int    `json:"frontendId"`
}

type ProgramSupplier struct {
	SupplierType      string `json:"supplierType"`
	AccountType       string `json:"accountType"`
	ProgramProviderId string `json:"programProviderId"`
	Name              string `json:"name"`
	Level             int    `json:"level"`
	PageUrl           string `json:"pageUrl"`
	Introduction      string `json:"introduction"`
}

type ProgramStatus string

const (
	ProgramStatus_ONAIR ProgramStatus = "ON_AIR"
	ProgramStatus_ENDED ProgramStatus = "ENDED"
)

type MediaServerType string

const (
	MediaServerType_DMC   MediaServerType = "DMC"
	MediaServerType_DLIVE MediaServerType = "DLIVE"
)

type Program struct {
	Status            ProgramStatus   `json:"status"`
	MediaServerType   MediaServerType `json:"mediaServerType"`
	Supplier          ProgramSupplier `json:"supplier"`
	NicoliveProgramId string          `json:"nicoliveProgramId"`
	Title             string          `json:"title"`
	Description       string          `json:"description"`
	WatchPageUrl      string          `json:"watchPageUrl"`

	OpenTime         int `json:"openTime"`
	BeginTime        int `json:"beginTime"`
	VposBaseTime     int `json:"vposBaseTime"`
	EndTime          int `json:"endTime"`
	ScheduledEndTime int `json:"scheduledEndTime"`

	IsPrivate                    bool `json:"isPrivate"`
	IsTest                       bool `json:"isTest"`
	IsFollowerOnly               bool `json:"isFollowerOnly"`
	IsNicoadEnabled              bool `json:"isNicoadEnabled"`
	IsGiftEnabled                bool `json:"isGiftEnabled"`
	IsChasePlayEnabled           bool `json:"isChasePlayEnabled"`
	IsTimeshiftDownloadEnabled   bool `json:"isTimeshiftDownloadEnabled"`
	IsPremiumAppealBannerEnabled bool `json:"isPremiumAppealBannerEnabled"`
	IsRecommendEnabled           bool `json:"isRecommendEnabled"`
	IsEmotionEnabled             bool `json:"isEmotionEnabled"`
}

type ProgramData struct {
	Site    Site    `json:"site"`
	Stream  Stream  `json:"stream"`
	Program Program `json:"program"`
}
