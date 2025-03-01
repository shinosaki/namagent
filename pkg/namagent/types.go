package namagent

import (
	"context"
	"net/http"
	"time"

	"github.com/shinosaki/namagent/internal/config"
)

type Alert func(
	config *config.Config,
	ctx context.Context,
	cancel context.CancelFunc,
) error

type Plugin struct {
	ExtractId func(input string) (id string)
	Session   func(config *config.Config) (client *http.Client)
	Client    func(
		id string,
		config *config.Config,
		client *http.Client,
		ctx context.Context,
	) (commentChan chan any, streamDataChan chan StreamData, err error)
}

type StreamData struct {
	URL       string
	Cookies   []*http.Cookie
	Output    string
	Extension string
	Template  Template
}

type Template struct {
	AuthorId   string
	AuthorName string

	ProgramId    string
	ProgramTitle string

	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
}
